// Deep 预构建智能体，支持子智能体、文件后端与深度任务迭代
package deep

import (
	"context"
	"fmt"

	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/adk/filesystem"
	"github.com/cloudwego/eino/adk/prebuilt/deep"
	"github.com/cloudwego/eino/components/model"
	"github.com/leitiannet/einolib"
)

type DeepAgentConfig struct {
	deep.Config // 内嵌deep.Config
}

func NewDeepAgentConfig(deepAgentOptions ...DeepAgentOption) *DeepAgentConfig {
	deepAgentConfig := &DeepAgentConfig{}
	deepAgentConfig.MaxIteration = 10 // 默认最大迭代次数为10
	einolib.ApplyOptions(deepAgentConfig, deepAgentOptions)
	return deepAgentConfig
}

type DeepAgentOption func(deepAgentConfig *DeepAgentConfig)

var (
	WithName                         = einolib.MakeOption(func(c *DeepAgentConfig, v string) { c.Name = v })
	WithDescription                  = einolib.MakeOption(func(c *DeepAgentConfig, v string) { c.Description = v })
	WithChatModel                    = einolib.MakeOption(func(c *DeepAgentConfig, v model.ToolCallingChatModel) { c.ChatModel = v })
	WithInstruction                  = einolib.MakeOption(func(c *DeepAgentConfig, v string) { c.Instruction = v })
	WithToolsConfig                  = einolib.MakeOption(func(c *DeepAgentConfig, v adk.ToolsConfig) { c.ToolsConfig = v })
	WithBackend                      = einolib.MakeOption(func(c *DeepAgentConfig, v filesystem.Backend) { c.Backend = v })
	WithShell                        = einolib.MakeOption(func(c *DeepAgentConfig, v filesystem.Shell) { c.Shell = v })
	WithStreamingShell               = einolib.MakeOption(func(c *DeepAgentConfig, v filesystem.StreamingShell) { c.StreamingShell = v })
	WithMaxIteration                 = einolib.MakeOption(func(c *DeepAgentConfig, v int) { c.MaxIteration = v })
	WithWithoutWriteTodos            = einolib.MakeOption(func(c *DeepAgentConfig, v bool) { c.WithoutWriteTodos = v })
	WithWithoutGeneralSubAgent       = einolib.MakeOption(func(c *DeepAgentConfig, v bool) { c.WithoutGeneralSubAgent = v })
	WithModelRetryConfig             = einolib.MakeOption(func(c *DeepAgentConfig, v *adk.ModelRetryConfig) { c.ModelRetryConfig = v })
	WithOutputKey                    = einolib.MakeOption(func(c *DeepAgentConfig, v string) { c.OutputKey = v })
	WithSubAgents                    = einolib.MakeAppendOption(func(c *DeepAgentConfig) *[]adk.Agent { return &c.SubAgents })
	WithMiddlewares                  = einolib.MakeAppendOption(func(c *DeepAgentConfig) *[]adk.AgentMiddleware { return &c.Middlewares })
	WithTaskToolDescriptionGenerator = einolib.MakeOption(func(c *DeepAgentConfig, v func(ctx context.Context, availableAgents []adk.Agent) (string, error)) {
		c.TaskToolDescriptionGenerator = v
	})
)

func NewDeepAgent(ctx context.Context, agentConfig *einolib.AgentConfig, deepAgentConfig *DeepAgentConfig) (adk.Agent, error) {
	agentConfig.ApplyNameAndDescription(&deepAgentConfig.Name, &deepAgentConfig.Description)
	if deepAgentConfig.ChatModel == nil {
		return nil, fmt.Errorf("deep agent: model is required")
	}
	for _, sa := range deepAgentConfig.SubAgents {
		if sa == nil {
			return nil, fmt.Errorf("deep agent: subAgent is nil")
		}
	}
	return deep.New(ctx, &deepAgentConfig.Config)
}

func createDeepAgent(ctx context.Context, agentConfig *einolib.AgentConfig, specificConfig interface{}) (adk.Agent, error) {
	deepAgentConfig, err := einolib.ParseSpecificConfig(specificConfig, func() *DeepAgentConfig { return NewDeepAgentConfig() })
	if err != nil {
		return nil, err
	}
	return NewDeepAgent(ctx, agentConfig, deepAgentConfig)
}

func init() {
	if err := einolib.RegisterAgentConstructFunc(einolib.AgentTypeDeep, einolib.AgentNameGeneral, createDeepAgent, (*DeepAgentConfig)(nil)); err != nil {
		einolib.GetLogger().Errorf("register agent %s failed: %v", einolib.AgentTypeDeep, err)
	}
}
