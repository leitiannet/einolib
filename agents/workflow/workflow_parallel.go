package workflow

import (
	"context"

	"github.com/cloudwego/eino/adk"
	"github.com/leitiannet/einolib"
)

type ParallelWorkflowAgentConfig struct {
	WorkflowAgentConfigCommon
}

func NewParallelWorkflowAgentConfig(workflowAgentOptions ...WorkflowAgentOption) *ParallelWorkflowAgentConfig {
	c := &ParallelWorkflowAgentConfig{}
	einolib.ApplyOptionsVariadic(&c.WorkflowAgentConfigCommon, workflowAgentOptions...)
	return c
}

func NewParallelWorkflowAgent(ctx context.Context, agentConfig *einolib.AgentConfig, cfg *ParallelWorkflowAgentConfig) (adk.Agent, error) {
	if err := validateAndApplyAgentMeta(agentConfig, &cfg.WorkflowAgentConfigCommon); err != nil {
		return nil, err
	}
	return adk.NewParallelAgent(ctx, &adk.ParallelAgentConfig{
		Name:        cfg.Name,
		Description: cfg.Description,
		SubAgents:   cfg.SubAgents,
	})
}

func createParallelAgent(ctx context.Context, agentConfig *einolib.AgentConfig, specificConfig interface{}) (adk.Agent, error) {
	cfg, err := einolib.ParseSpecificConfig(specificConfig, func() *ParallelWorkflowAgentConfig { return NewParallelWorkflowAgentConfig() })
	if err != nil {
		return nil, err
	}
	return NewParallelWorkflowAgent(ctx, agentConfig, cfg)
}

func init() {
	if err := einolib.RegisterAgentConstructFunc(einolib.AgentTypeWorkflowParallel, einolib.AgentNameGeneral, createParallelAgent, (*ParallelWorkflowAgentConfig)(nil)); err != nil {
		einolib.GetLogger().Errorf("register agent %s failed: %v", einolib.AgentTypeWorkflowParallel, err)
	}
}
