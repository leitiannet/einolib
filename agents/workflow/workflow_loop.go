package workflow

import (
	"context"

	"github.com/cloudwego/eino/adk"
	"github.com/leitiannet/einolib"
)

const AgentTypeWorkflowLoop einolib.AgentType = "workflow_loop"

type LoopWorkflowAgentConfig struct {
	WorkflowAgentConfigCommon
	MaxIterations int // 最大迭代次数
}

// 可同时传入WorkflowAgentOption与LoopWorkflowAgentOption
func NewLoopWorkflowAgentConfig(options ...interface{}) *LoopWorkflowAgentConfig {
	c := &LoopWorkflowAgentConfig{MaxIterations: 10}
	for _, opt := range options {
		if opt == nil {
			continue
		}
		switch o := opt.(type) {
		case LoopWorkflowAgentOption:
			o(c)
		case WorkflowAgentOption:
			o(&c.WorkflowAgentConfigCommon)
		case func(*LoopWorkflowAgentConfig):
			o(c)
		case func(*WorkflowAgentConfigCommon):
			o(&c.WorkflowAgentConfigCommon)
		default:
			einolib.GetLogger().Warnf("workflow.NewLoopWorkflowAgentConfig: skip unsupported option type %T", opt)
		}
	}
	return c
}

type LoopWorkflowAgentOption func(*LoopWorkflowAgentConfig)

var (
	WithMaxIterations = einolib.MakeOption(func(c *LoopWorkflowAgentConfig, v int) { c.MaxIterations = v })
)

func NewLoopWorkflowAgent(ctx context.Context, agentConfig *einolib.AgentConfig, cfg *LoopWorkflowAgentConfig) (adk.Agent, error) {
	if err := validateAndApplyAgentMeta(agentConfig, &cfg.WorkflowAgentConfigCommon); err != nil {
		return nil, err
	}
	return adk.NewLoopAgent(ctx, &adk.LoopAgentConfig{
		Name:          cfg.Name,
		Description:   cfg.Description,
		SubAgents:     cfg.SubAgents,
		MaxIterations: cfg.MaxIterations,
	})
}

func createLoopAgent(ctx context.Context, agentConfig *einolib.AgentConfig, specificConfig interface{}) (adk.Agent, error) {
	cfg, err := einolib.ParseSpecificConfig(specificConfig, func() *LoopWorkflowAgentConfig { return NewLoopWorkflowAgentConfig() })
	if err != nil {
		return nil, err
	}
	return NewLoopWorkflowAgent(ctx, agentConfig, cfg)
}

func init() {
	if err := einolib.RegisterAgentConstructFunc(AgentTypeWorkflowLoop, einolib.GeneralAgentName, createLoopAgent, (*LoopWorkflowAgentConfig)(nil)); err != nil {
		einolib.GetLogger().Errorf("register agent %s failed: %v", AgentTypeWorkflowLoop, err)
	}
}
