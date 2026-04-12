package workflow

import (
	"context"

	"github.com/cloudwego/eino/adk"
	"github.com/leitiannet/einolib"
)

const AgentTypeWorkflowParallel einolib.AgentType = "workflow_parallel"

type ParallelWorkflowAgentConfig struct {
	WorkflowAgentConfigCommon
}

func NewParallelWorkflowAgentConfig(workflowAgentOptions ...WorkflowAgentOption) *ParallelWorkflowAgentConfig {
	c := &ParallelWorkflowAgentConfig{}
	einolib.ApplyOptionsVariadic(&c.WorkflowAgentConfigCommon, workflowAgentOptions...)
	return c
}

func NewParallelAgent(ctx context.Context, agentConfig *einolib.AgentConfig, specificConfig interface{}) (adk.Agent, error) {
	cfg, err := einolib.ParseSpecificConfig(specificConfig, func() *ParallelWorkflowAgentConfig { return NewParallelWorkflowAgentConfig() })
	if err != nil {
		return nil, err
	}
	if err := validateAndApplyAgentMeta(agentConfig, &cfg.WorkflowAgentConfigCommon); err != nil {
		return nil, err
	}
	return adk.NewParallelAgent(ctx, &adk.ParallelAgentConfig{
		Name:        cfg.Name,
		Description: cfg.Description,
		SubAgents:   cfg.SubAgents,
	})
}

func init() {
	if err := einolib.RegisterAgentConstructFunc(AgentTypeWorkflowParallel, einolib.GeneralAgentName, NewParallelAgent, (*ParallelWorkflowAgentConfig)(nil)); err != nil {
		einolib.GetLogger().Errorf("register agent %s failed: %v", AgentTypeWorkflowParallel, err)
	}
}
