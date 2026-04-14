// 火山引擎AgentKit沙箱执行器，支持文件操作和命令执行
package agentkit

import (
	"context"
	"net/http"

	einoagentkit "github.com/cloudwego/eino-ext/adk/backend/agentkit"
	"github.com/leitiannet/einolib"
)

const (
	OperatorTypeAgentKit einolib.OperatorType = "agentkit"
)

type AgentKitOperatorConfig struct {
	einoagentkit.Config // 内嵌结构体
}

func NewAgentKitOperatorConfig(agentKitOperatorOptions ...AgentKitOperatorOption) *AgentKitOperatorConfig {
	agentKitOperatorConfig := &AgentKitOperatorConfig{}
	einolib.ApplyOptions(agentKitOperatorConfig, agentKitOperatorOptions)
	return agentKitOperatorConfig
}

type AgentKitOperatorOption func(*AgentKitOperatorConfig)

var (
	WithAccessKeyID      = einolib.MakeOption(func(c *AgentKitOperatorConfig, v string) { c.AccessKeyID = v })
	WithSecretAccessKey  = einolib.MakeOption(func(c *AgentKitOperatorConfig, v string) { c.SecretAccessKey = v })
	WithHTTPClient       = einolib.MakeOption(func(c *AgentKitOperatorConfig, v *http.Client) { c.HTTPClient = v })
	WithRegion           = einolib.MakeOption(func(c *AgentKitOperatorConfig, v einoagentkit.Region) { c.Region = v })
	WithToolID           = einolib.MakeOption(func(c *AgentKitOperatorConfig, v string) { c.ToolID = v })
	WithSessionID        = einolib.MakeOption(func(c *AgentKitOperatorConfig, v string) { c.SessionID = v })
	WithUserSessionID    = einolib.MakeOption(func(c *AgentKitOperatorConfig, v string) { c.UserSessionID = v })
	WithSessionTTL       = einolib.MakeOption(func(c *AgentKitOperatorConfig, v int) { c.SessionTTL = v })
	WithExecutionTimeout = einolib.MakeOption(func(c *AgentKitOperatorConfig, v int) { c.ExecutionTimeout = v })
)

func NewAgentKitOperator(ctx context.Context, agentKitOperatorConfig *AgentKitOperatorConfig) (einolib.Operator, error) {
	sandbox, err := einoagentkit.NewSandboxToolBackend(&agentKitOperatorConfig.Config)
	if err != nil {
		return nil, err
	}
	return &AgentKitOperator{SandboxTool: sandbox}, nil
}

func createAgentKitOperator(ctx context.Context, operatorConfig *einolib.OperatorConfig, specificConfig interface{}) (einolib.Operator, error) {
	agentKitOperatorConfig, err := einolib.ParseSpecificConfig(specificConfig, func() *AgentKitOperatorConfig { return NewAgentKitOperatorConfig() })
	if err != nil {
		return nil, err
	}
	return NewAgentKitOperator(ctx, agentKitOperatorConfig)
}

func init() {
	if err := einolib.RegisterOperatorConstructFunc(OperatorTypeAgentKit, einolib.GeneralOperatorName, createAgentKitOperator, (*AgentKitOperatorConfig)(nil)); err != nil {
		einolib.GetLogger().Errorf("register operator %s failed: %v", OperatorTypeAgentKit, err)
	}
}
