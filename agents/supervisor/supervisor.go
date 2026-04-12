package supervisor

import (
	"context"
	"fmt"

	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/adk/prebuilt/supervisor"
	"github.com/leitiannet/einolib"
)

const (
	AgentTypeSupervisor einolib.AgentType = "supervisor"
)

type SupervisorAgentConfig struct {
	supervisor.Config // 内嵌supervisor.Config
}

func NewSupervisorAgentConfig(supervisorAgentOptions ...SupervisorAgentOption) *SupervisorAgentConfig {
	supervisorAgentConfig := &SupervisorAgentConfig{}
	einolib.ApplyOptions(supervisorAgentConfig, supervisorAgentOptions)
	return supervisorAgentConfig
}

type SupervisorAgentOption func(supervisorAgentConfig *SupervisorAgentConfig)

var (
	WithSupervisor = einolib.MakeOption(func(c *SupervisorAgentConfig, v adk.Agent) { c.Supervisor = v })
	WithSubAgents  = einolib.MakeAppendOption(func(c *SupervisorAgentConfig) *[]adk.Agent { return &c.SubAgents })
)

func NewSupervisorAgent(ctx context.Context, agentConfig *einolib.AgentConfig, specificConfig interface{}) (adk.Agent, error) {
	supervisorAgentConfig, err := einolib.ParseSpecificConfig(specificConfig, func() *SupervisorAgentConfig { return NewSupervisorAgentConfig() })
	if err != nil {
		return nil, err
	}
	if supervisorAgentConfig.Supervisor == nil {
		return nil, fmt.Errorf("supervisor agent: supervisor is required")
	}
	if len(supervisorAgentConfig.SubAgents) == 0 {
		return nil, fmt.Errorf("supervisor agent: subAgents is required")
	}
	for _, subAgent := range supervisorAgentConfig.SubAgents {
		if subAgent == nil {
			return nil, fmt.Errorf("supervisor agent: subAgent is nil")
		}
	}
	return supervisor.New(ctx, &supervisorAgentConfig.Config)
}

func init() {
	if err := einolib.RegisterAgentConstructFunc(AgentTypeSupervisor, einolib.GeneralAgentName, NewSupervisorAgent, (*SupervisorAgentConfig)(nil)); err != nil {
		einolib.GetLogger().Errorf("register agent %s failed: %v", AgentTypeSupervisor, err)
	}
}
