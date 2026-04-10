package einolib

import (
	"github.com/cloudwego/eino/schema"
)

// 简化调用方式，不需要记住多个函数
// 减少噪声参数，不需要的参数不传入

// 模板项接口，统一消息和占位符的构建入口
type TemplateItem interface {
	templateItems() []schema.MessagesTemplate
}

// 占位符配置
type PlaceholderConfig struct {
	Key      string // 占位符键
	Optional bool   // 是否可选
}

func (p *PlaceholderConfig) templateItems() []schema.MessagesTemplate {
	if p == nil || p.Key == "" {
		return nil
	}
	return []schema.MessagesTemplate{schema.MessagesPlaceholder(p.Key, p.Optional)}
}

// 消息配置
type MessageConfig struct {
	schema.Message // 内嵌 schema.Message（会话数据）
}

func (m *MessageConfig) templateItems() []schema.MessagesTemplate {
	if msg := m.toMessage(); msg != nil {
		return []schema.MessagesTemplate{msg}
	}
	return nil
}

func (m *MessageConfig) toMessage() *schema.Message {
	if m == nil {
		return nil
	}
	msg := m.Message
	return &msg
}

// 消息选项
type MsgOption func(*MessageConfig)

var (
	// 设置工具调用（Assistant 角色）
	WithToolCalls = MakeOption(func(c *MessageConfig, v []schema.ToolCall) { c.ToolCalls = v })
	// 设置工具调用 ID（Tool 角色）
	WithToolCallID = MakeOption(func(c *MessageConfig, v string) { c.ToolCallID = v })
	// 设置工具名称（Tool 角色）
	WithMsgToolName = MakeOption(func(c *MessageConfig, v string) { c.ToolName = v })
	// 设置自定义扩展参数
	WithExtra = MakeOption(func(c *MessageConfig, v map[string]any) { c.Extra = v })
)

// 快捷构造消息配置
// system：系统指令，通常放在 messages 最前面
// user：用户输入
// assistant：模型回复
// tool：工具调用结果
func Msg(role schema.RoleType, content string, opts ...MsgOption) *MessageConfig {
	if role == "" {
		return nil
	}
	m := &MessageConfig{Message: schema.Message{Role: role, Content: content}}
	ApplyOptions(m, opts)
	return m
}

// 快捷构造占位符配置
func Placeholder(key string, optional ...bool) *PlaceholderConfig {
	if key == "" {
		return nil
	}
	opt := false
	if len(optional) > 0 {
		opt = optional[0]
	}
	return &PlaceholderConfig{Key: key, Optional: opt}
}

// 创建消息列表
func NewMessages(configs ...*MessageConfig) []*schema.Message {
	messages := make([]*schema.Message, 0, len(configs))
	for _, cfg := range configs {
		if cfg == nil {
			continue
		}
		messages = append(messages, cfg.toMessage())
	}
	return messages
}

// 创建消息模板列表，同时接收 MessageConfig 和 PlaceholderConfig
func NewTemplates(items ...TemplateItem) []schema.MessagesTemplate {
	templates := make([]schema.MessagesTemplate, 0, len(items))
	for _, item := range items {
		if item == nil {
			continue
		}
		templates = append(templates, item.templateItems()...)
	}
	return templates
}