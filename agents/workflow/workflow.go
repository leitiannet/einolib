// Workflow 系列预构建智能体，按顺序、并行或循环编排子智能体
package workflow

import (
	"fmt"

	"github.com/cloudwego/eino/adk"
	"github.com/leitiannet/einolib"
)

// 共用配置
type WorkflowAgentConfigCommon struct {
	Name        string      // 智能体名称
	Description string      // 智能体描述
	SubAgents   []adk.Agent // 子智能体列表
}

type WorkflowAgentOption func(*WorkflowAgentConfigCommon)

var (
	WithName        = einolib.MakeOption(func(c *WorkflowAgentConfigCommon, v string) { c.Name = v })
	WithDescription = einolib.MakeOption(func(c *WorkflowAgentConfigCommon, v string) { c.Description = v })
	WithSubAgents   = einolib.MakeAppendOption(func(c *WorkflowAgentConfigCommon) *[]adk.Agent { return &c.SubAgents })
)

func validateAndApplyAgentMeta(agentConfig *einolib.AgentConfig, common *WorkflowAgentConfigCommon) error {
	agentConfig.ApplyNameAndDescription(&common.Name, &common.Description)
	if common.Name == "" {
		return fmt.Errorf("workflow agent: name is required")
	}
	if common.Description == "" {
		return fmt.Errorf("workflow agent: description is required")
	}
	if len(common.SubAgents) == 0 {
		return fmt.Errorf("workflow agent: subAgents is required")
	}
	for _, subAgent := range common.SubAgents {
		if subAgent == nil {
			return fmt.Errorf("workflow agent: subAgent is nil")
		}
	}
	return nil
}
