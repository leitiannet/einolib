package workflow

import (
	"context"

	"github.com/cloudwego/eino/adk"
	"github.com/leitiannet/einolib"
)

func NewSequentialAgent(ctx context.Context, agentConfig *einolib.AgentConfig, specificConfig interface{}) (adk.Agent, error) {
	workflowAgentConfig, err := resolveWorkflowAgentConfig(agentConfig, specificConfig)
	if err != nil {
		return nil, err
	}
	return adk.NewSequentialAgent(ctx, &adk.SequentialAgentConfig{
		Name:        workflowAgentConfig.Name,
		Description: workflowAgentConfig.Description,
		SubAgents:   workflowAgentConfig.SubAgents,
	})
}

func init() {
	if err := einolib.RegisterAgentConstructFunc(AgentTypeSequential, einolib.GeneralAgentName, NewSequentialAgent); err != nil {
		einolib.GetLogger().Errorf("register agent %s failed: %v", AgentTypeSequential, err)
	}
}
