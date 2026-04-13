// 将Agents.md内容注入到模型输入消息中
package agentsmd

import (
	"context"

	"github.com/cloudwego/eino/adk"
	amdmiddleware "github.com/cloudwego/eino/adk/middlewares/agentsmd"
	"github.com/leitiannet/einolib"
)

const (
	MiddlewareTypeAgentsMD einolib.MiddlewareType = "agentsmd"
)

type AgentsMDMiddlewareConfig struct {
	amdmiddleware.Config // 内嵌结构体
}

func NewAgentsMDMiddlewareConfig(agentsMDMiddlewareOptions ...AgentsMDMiddlewareOption) *AgentsMDMiddlewareConfig {
	config := &AgentsMDMiddlewareConfig{}
	einolib.ApplyOptions(config, agentsMDMiddlewareOptions)
	return config
}

type AgentsMDMiddlewareOption func(*AgentsMDMiddlewareConfig)

var (
	WithBackend             = einolib.MakeOption(func(c *AgentsMDMiddlewareConfig, v amdmiddleware.Backend) { c.Backend = v })
	WithAgentsMDFiles       = einolib.MakeAppendOption(func(c *AgentsMDMiddlewareConfig) *[]string { return &c.AgentsMDFiles })
	WithAllAgentsMDMaxBytes = einolib.MakeOption(func(c *AgentsMDMiddlewareConfig, v int) { c.AllAgentsMDMaxBytes = v })
	WithOnLoadWarning       = einolib.MakeOption(func(c *AgentsMDMiddlewareConfig, v func(string, error)) { c.OnLoadWarning = v })
)

func NewAgentsMDMiddleware(ctx context.Context, config *AgentsMDMiddlewareConfig) (adk.ChatModelAgentMiddleware, error) {
	return amdmiddleware.New(ctx, &config.Config)
}

func createMiddleware(ctx context.Context, middlewareConfig *einolib.MiddlewareConfig, specificConfig interface{}) (adk.ChatModelAgentMiddleware, error) {
	agentsMDMiddlewareConfig, err := einolib.ParseSpecificConfig(specificConfig, func() *AgentsMDMiddlewareConfig { return NewAgentsMDMiddlewareConfig() })
	if err != nil {
		return nil, err
	}
	return NewAgentsMDMiddleware(ctx, agentsMDMiddlewareConfig)
}

func init() {
	if err := einolib.RegisterMiddlewareConstructFunc(MiddlewareTypeAgentsMD, einolib.GeneralMiddlewareName, createMiddleware, (*AgentsMDMiddlewareConfig)(nil)); err != nil {
		einolib.GetLogger().Errorf("register middleware %s failed: %v", MiddlewareTypeAgentsMD, err)
	}
}
