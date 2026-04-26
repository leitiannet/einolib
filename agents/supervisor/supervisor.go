// Supervisor 预构建智能体，由监督智能体协调多个子智能体
package supervisor

import (
	"context"
	"fmt"

	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/adk/prebuilt/supervisor"
	"github.com/leitiannet/einolib"
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

func NewSupervisorAgent(ctx context.Context, agentConfig *einolib.AgentConfig, supervisorAgentConfig *SupervisorAgentConfig) (adk.Agent, error) {
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

func createSupervisorAgent(ctx context.Context, agentConfig *einolib.AgentConfig, specificConfig interface{}) (adk.Agent, error) {
	supervisorAgentConfig, err := einolib.ParseSpecificConfig(specificConfig, func() *SupervisorAgentConfig { return NewSupervisorAgentConfig() })
	if err != nil {
		return nil, err
	}
	return NewSupervisorAgent(ctx, agentConfig, supervisorAgentConfig)
}

func init() {
	if err := einolib.RegisterAgentConstructFunc(einolib.AgentTypeSupervisor, einolib.AgentNameGeneral, createSupervisorAgent, (*SupervisorAgentConfig)(nil)); err != nil {
		einolib.GetLogger().Errorf("register agent %s failed: %v", einolib.AgentTypeSupervisor, err)
	}
}
