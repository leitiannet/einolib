package einolib

import (
	"context"
	"fmt"

	"github.com/cloudwego/eino/adk"
)

// 运行器配置
type RunnerConfig struct {
	adk.RunnerConfig               // 内嵌adk.RunnerConfig
	AgentOptions     []AgentOption // 通过AgentOption构建Agent（Agent为nil时生效）
}

func NewRunnerConfig(runnerOptions ...RunnerOption) *RunnerConfig {
	runnerConfig := &RunnerConfig{}
	ApplyOptions(runnerConfig, runnerOptions)
	return runnerConfig
}

// 运行器选项
type RunnerOption func(runnerConfig *RunnerConfig)

var (
	WithAgent           = MakeOption(func(c *RunnerConfig, v adk.Agent) { c.Agent = v })
	WithEnableStreaming  = MakeOption(func(c *RunnerConfig, v bool) { c.EnableStreaming = v })
	WithCheckPointStore = MakeOption(func(c *RunnerConfig, v adk.CheckPointStore) { c.CheckPointStore = v })
	WithAgentOptions    = MakeAppendOption(func(c *RunnerConfig) *[]AgentOption { return &c.AgentOptions })
)

func NewRunner(ctx context.Context, runnerOptions ...RunnerOption) (*adk.Runner, error) {
	runnerConfig := NewRunnerConfig(runnerOptions...)
	if runnerConfig.Agent == nil {
		var err error
		runnerConfig.Agent, err = NewAgent(ctx, runnerConfig.AgentOptions...)
		if err != nil {
			return nil, err
		}
	}
	if runnerConfig.Agent == nil {
		return nil, fmt.Errorf("runner: agent is nil")
	}
	return adk.NewRunner(ctx, runnerConfig.RunnerConfig), nil
}

func MustNewRunner(ctx context.Context, runnerOptions ...RunnerOption) *adk.Runner {
	runner, err := NewRunner(ctx, runnerOptions...)
	if err != nil {
		panic(err)
	}
	if runner == nil {
		panic("MustNewRunner failed: instance is nil")
	}
	return runner
}
