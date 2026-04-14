// ChatModelAgent是Eino ADK中一个核心预构建的Agent，它封装了与大语言模型（LLM）进行交互、并支持使用工具来完成任务的复杂逻辑，内置ReAct模式
package chatmodel

import (
	"context"
	"fmt"

	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/components/tool"
	"github.com/leitiannet/einolib"
)

const (
	AgentTypeChatModel einolib.AgentType = "chatmodel"
)

type ChatModelAgentConfig struct {
	adk.ChatModelAgentConfig // 内嵌adk.ChatModelAgentConfig
}

func NewChatModelAgentConfig(chatModelAgentOptions ...ChatModelAgentOption) *ChatModelAgentConfig {
	chatModelAgentConfig := &ChatModelAgentConfig{}
	einolib.ApplyOptions(chatModelAgentConfig, chatModelAgentOptions)
	return chatModelAgentConfig
}

type ChatModelAgentOption func(chatModelAgentConfig *ChatModelAgentConfig)

var (
	WithName             = einolib.MakeOption(func(c *ChatModelAgentConfig, v string) { c.Name = v })
	WithDescription      = einolib.MakeOption(func(c *ChatModelAgentConfig, v string) { c.Description = v })
	WithInstruction      = einolib.MakeOption(func(c *ChatModelAgentConfig, v string) { c.Instruction = v })
	WithModel            = einolib.MakeOption(func(c *ChatModelAgentConfig, v model.ToolCallingChatModel) { c.Model = v })
	WithToolsConfig      = einolib.MakeOption(func(c *ChatModelAgentConfig, v adk.ToolsConfig) { c.ToolsConfig = v })
	WithExit             = einolib.MakeOption(func(c *ChatModelAgentConfig, v tool.BaseTool) { c.Exit = v })
	WithOutputKey        = einolib.MakeOption(func(c *ChatModelAgentConfig, v string) { c.OutputKey = v })
	WithMaxIterations    = einolib.MakeOption(func(c *ChatModelAgentConfig, v int) { c.MaxIterations = v })
	WithMiddlewares      = einolib.MakeAppendOption(func(c *ChatModelAgentConfig) *[]adk.AgentMiddleware { return &c.Middlewares })
	WithHandlers         = einolib.MakeAppendOption(func(c *ChatModelAgentConfig) *[]adk.ChatModelAgentMiddleware { return &c.Handlers })
	WithModelRetryConfig = einolib.MakeOption(func(c *ChatModelAgentConfig, v *adk.ModelRetryConfig) { c.ModelRetryConfig = v })
)

func NewChatModelAgent(ctx context.Context, agentConfig *einolib.AgentConfig, chatModelAgentConfig *ChatModelAgentConfig) (adk.Agent, error) {
	agentConfig.ApplyNameAndDescription(&chatModelAgentConfig.Name, &chatModelAgentConfig.Description)
	if chatModelAgentConfig.Model == nil {
		return nil, fmt.Errorf("chatmodel agent: model is required")
	}
	return adk.NewChatModelAgent(ctx, &chatModelAgentConfig.ChatModelAgentConfig)
}

func createChatModelAgent(ctx context.Context, agentConfig *einolib.AgentConfig, specificConfig interface{}) (adk.Agent, error) {
	chatModelAgentConfig, err := einolib.ParseSpecificConfig(specificConfig, func() *ChatModelAgentConfig { return NewChatModelAgentConfig() })
	if err != nil {
		return nil, err
	}
	return NewChatModelAgent(ctx, agentConfig, chatModelAgentConfig)
}

func init() {
	if err := einolib.RegisterAgentConstructFunc(AgentTypeChatModel, einolib.GeneralAgentName, createChatModelAgent, (*ChatModelAgentConfig)(nil)); err != nil {
		einolib.GetLogger().Errorf("register agent %s failed: %v", AgentTypeChatModel, err)
	}
}
