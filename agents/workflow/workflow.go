package workflow

import (
	"fmt"

	"github.com/cloudwego/eino/adk"
	"github.com/leitiannet/einolib"
)

const (
	AgentTypeSequential einolib.AgentType = "workflow_sequential"
	AgentTypeParallel   einolib.AgentType = "workflow_parallel"
	AgentTypeLoop       einolib.AgentType = "workflow_loop"
)

type WorkflowAgentConfig struct {
	Name          string      // 智能体名称
	Description   string      // 智能体描述
	SubAgents     []adk.Agent // 子智能体列表
	MaxIterations int         // 最大循环次数（仅 LoopAgent 使用）
}

func NewWorkflowAgentConfig(workflowAgentOptions ...WorkflowAgentOption) *WorkflowAgentConfig {
	workflowAgentConfig := &WorkflowAgentConfig{}
	workflowAgentConfig.MaxIterations = 10 // 默认最大循环次数为10
	einolib.ApplyOptions(workflowAgentConfig, workflowAgentOptions)
	return workflowAgentConfig
}

type WorkflowAgentOption func(workflowAgentConfig *WorkflowAgentConfig)

var (
	WithName          = einolib.MakeOption(func(c *WorkflowAgentConfig, v string) { c.Name = v })
	WithDescription   = einolib.MakeOption(func(c *WorkflowAgentConfig, v string) { c.Description = v })
	WithMaxIterations = einolib.MakeOption(func(c *WorkflowAgentConfig, v int) { c.MaxIterations = v })
	WithSubAgents     = einolib.MakeAppendOption(func(c *WorkflowAgentConfig) *[]adk.Agent { return &c.SubAgents })
)

func resolveWorkflowAgentConfig(agentConfig *einolib.AgentConfig, specificConfig interface{}) (*WorkflowAgentConfig, error) {
	workflowAgentConfig, err := einolib.ParseSpecificConfig(specificConfig, func() *WorkflowAgentConfig { return NewWorkflowAgentConfig() })
	if err != nil {
		return nil, err
	}
	if len(workflowAgentConfig.SubAgents) == 0 {
		return nil, fmt.Errorf("workflow agent: subAgents is required")
	}
	for _, subAgent := range workflowAgentConfig.SubAgents {
		if subAgent == nil {
			return nil, fmt.Errorf("workflow agent: subAgent is nil")
		}
	}
	agentConfig.ApplyNameAndDescription(&workflowAgentConfig.Name, &workflowAgentConfig.Description)
	return workflowAgentConfig, nil
}
