package workflow

import (
	"context"

	"github.com/cloudwego/eino/adk"
	"github.com/leitiannet/einolib"
)

func NewLoopAgent(ctx context.Context, agentConfig *einolib.AgentConfig, specificConfig interface{}) (adk.Agent, error) {
	workflowAgentConfig, err := resolveWorkflowAgentConfig(agentConfig, specificConfig)
	if err != nil {
		return nil, err
	}
	return adk.NewLoopAgent(ctx, &adk.LoopAgentConfig{
		Name:          workflowAgentConfig.Name,
		Description:   workflowAgentConfig.Description,
		SubAgents:     workflowAgentConfig.SubAgents,
		MaxIterations: workflowAgentConfig.MaxIterations,
	})
}

func init() {
	if err := einolib.RegisterAgentConstructFunc(AgentTypeLoop, einolib.GeneralAgentName, NewLoopAgent); err != nil {
		einolib.GetLogger().Errorf("register agent %s failed: %v", AgentTypeLoop, err)
	}
}
