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
	WithName          = einolib.MakeOption(func(c *ChatModelAgentConfig, v string) { c.Name = v })
	WithDescription   = einolib.MakeOption(func(c *ChatModelAgentConfig, v string) { c.Description = v })
	WithInstruction   = einolib.MakeOption(func(c *ChatModelAgentConfig, v string) { c.Instruction = v })
	WithModel         = einolib.MakeOption(func(c *ChatModelAgentConfig, v model.ToolCallingChatModel) { c.Model = v })
	WithToolsConfig   = einolib.MakeOption(func(c *ChatModelAgentConfig, v adk.ToolsConfig) { c.ToolsConfig = v })
	WithExit          = einolib.MakeOption(func(c *ChatModelAgentConfig, v tool.BaseTool) { c.Exit = v })
	WithOutputKey     = einolib.MakeOption(func(c *ChatModelAgentConfig, v string) { c.OutputKey = v })
	WithMaxIterations = einolib.MakeOption(func(c *ChatModelAgentConfig, v int) { c.MaxIterations = v })
)

func NewChatModelAgent(ctx context.Context, agentConfig *einolib.AgentConfig, specificConfig interface{}) (adk.Agent, error) {
	chatModelAgentConfig, err := einolib.ParseSpecificConfig(specificConfig, func() *ChatModelAgentConfig { return NewChatModelAgentConfig() })
	if err != nil {
		return nil, err
	}
	agentConfig.ApplyNameAndDescription(&chatModelAgentConfig.Name, &chatModelAgentConfig.Description)
	if chatModelAgentConfig.Model == nil {
		return nil, fmt.Errorf("chatmodel agent: model is required")
	}
	return adk.NewChatModelAgent(ctx, &chatModelAgentConfig.ChatModelAgentConfig)
}

func init() {
	if err := einolib.RegisterAgentConstructFunc(AgentTypeChatModel, einolib.GeneralAgentName, NewChatModelAgent, (*ChatModelAgentConfig)(nil)); err != nil {
		einolib.GetLogger().Errorf("register agent %s failed: %v", AgentTypeChatModel, err)
	}
}
