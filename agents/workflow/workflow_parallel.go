package workflow

import (
	"context"

	"github.com/cloudwego/eino/adk"
	"github.com/leitiannet/einolib"
)

func NewParallelAgent(ctx context.Context, agentConfig *einolib.AgentConfig, specificConfig interface{}) (adk.Agent, error) {
	workflowAgentConfig, err := resolveWorkflowAgentConfig(agentConfig, specificConfig)
	if err != nil {
		return nil, err
	}
	return adk.NewParallelAgent(ctx, &adk.ParallelAgentConfig{
		Name:        workflowAgentConfig.Name,
		Description: workflowAgentConfig.Description,
		SubAgents:   workflowAgentConfig.SubAgents,
	})
}

func init() {
	if err := einolib.RegisterAgentConstructFunc(AgentTypeParallel, einolib.GeneralAgentName, NewParallelAgent); err != nil {
		einolib.GetLogger().Errorf("register agent %s failed: %v", AgentTypeParallel, err)
	}
}
