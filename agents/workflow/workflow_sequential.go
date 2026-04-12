package workflow

import (
	"context"

	"github.com/cloudwego/eino/adk"
	"github.com/leitiannet/einolib"
)

const AgentTypeWorkflowSequential einolib.AgentType = "workflow_sequential"

type SequentialWorkflowAgentConfig struct {
	WorkflowAgentConfigCommon
}

func NewSequentialWorkflowAgentConfig(workflowAgentOptions ...WorkflowAgentOption) *SequentialWorkflowAgentConfig {
	c := &SequentialWorkflowAgentConfig{}
	einolib.ApplyOptionsVariadic(&c.WorkflowAgentConfigCommon, workflowAgentOptions...)
	return c
}

func NewSequentialAgent(ctx context.Context, agentConfig *einolib.AgentConfig, specificConfig interface{}) (adk.Agent, error) {
	cfg, err := einolib.ParseSpecificConfig(specificConfig, func() *SequentialWorkflowAgentConfig { return NewSequentialWorkflowAgentConfig() })
	if err != nil {
		return nil, err
	}
	if err := validateAndApplyAgentMeta(agentConfig, &cfg.WorkflowAgentConfigCommon); err != nil {
		return nil, err
	}
	return adk.NewSequentialAgent(ctx, &adk.SequentialAgentConfig{
		Name:        cfg.Name,
		Description: cfg.Description,
		SubAgents:   cfg.SubAgents,
	})
}

func init() {
	if err := einolib.RegisterAgentConstructFunc(AgentTypeWorkflowSequential, einolib.GeneralAgentName, NewSequentialAgent, (*SequentialWorkflowAgentConfig)(nil)); err != nil {
		einolib.GetLogger().Errorf("register agent %s failed: %v", AgentTypeWorkflowSequential, err)
	}
}
