package einolib

import "github.com/cloudwego/eino/schema"

type MessageConfig struct {
	RoleType           schema.RoleType
	Content            string
	ToolCalls          []schema.ToolCall
	ToolCallID         string
	ToolMessageOptions []schema.ToolMessageOption
}

func NewMessages(messageConfigs ...*MessageConfig) []*schema.Message {
	messages := make([]*schema.Message, 0, len(messageConfigs))
	for _, messageConfig := range messageConfigs {
		if messageConfig == nil {
			continue
		}
		switch messageConfig.RoleType {
		case schema.Assistant:
			messages = append(messages, schema.AssistantMessage(messageConfig.Content, messageConfig.ToolCalls))
		case schema.User:
			messages = append(messages, schema.UserMessage(messageConfig.Content))
		case schema.System:
			messages = append(messages, schema.SystemMessage(messageConfig.Content))
		case schema.Tool:
			messages = append(messages, schema.ToolMessage(messageConfig.Content, messageConfig.ToolCallID, messageConfig.ToolMessageOptions...))
		}
	}
	return messages
}

type MessagesTemplateConfig struct {
	MessageConfig
	PlaceholderKey      string
	PlaceholderOptional bool
}

func NewMessagesTemplates(messagesTemplateConfigs ...*MessagesTemplateConfig) []schema.MessagesTemplate {
	messagesTemplates := make([]schema.MessagesTemplate, 0, len(messagesTemplateConfigs))
	return messagesTemplates
}
