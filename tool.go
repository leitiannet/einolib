package einolib

import (
	"context"

	"github.com/cloudwego/eino/components"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/schema"
)

// 工具类型
type ToolType string

const (
	ToolTypeCustom  ToolType = "custom"  // 用户定义工具
	ToolTypeBuiltin ToolType = "builtin" // 官方内置工具
	ToolTypeMCP     ToolType = "mcp"     // MCP协议工具
)

const ToolNameGeneral = ComponentNameGeneral // 通用工具名称

// 工具描述
type ToolDescriber struct {
	componentDescriber
}

func NewToolDescriber(toolType ToolType, toolName string) *ToolDescriber {
	return &ToolDescriber{
		componentDescriber: NewComponentDescriber(components.ComponentOfTool, string(toolType), toolName),
	}
}

// 工具配置
type ToolConfig struct {
	componentConfig
}

func NewToolConfig(toolOptions ...ToolOption) *ToolConfig {
	toolConfig := &ToolConfig{
		componentConfig: NewComponentConfig(components.ComponentOfTool, "", ""),
	}
	ApplyOptions(toolConfig, toolOptions)
	return toolConfig
}

// 工具选项
type ToolOption func(toolConfig *ToolConfig)

var (
	WithToolType            = MakeOption(func(c *ToolConfig, v ToolType) { c.Type = string(v) })
	WithToolName            = MakeOption(func(c *ToolConfig, v string) { c.Name = v })
	WithToolDescription     = MakeOption(func(c *ToolConfig, v string) { c.Description = v })
	WithToolComponentConfig = MakeOption(func(c *ToolConfig, value interface{}) {
		desc, err := LookupToolDescriber(value)
		if err != nil {
			logger.Warnf("LookupToolDescriber failed: %v", err)
			return
		}
		if desc == nil {
			logger.Warnf("describer is nil for type %T", value)
			return
		}
		c.SetConfig(desc, value)
	})
)

// 工具构造器接口
type ToolConstructor interface {
	Construct(ctx context.Context, toolConfig *ToolConfig) ([]tool.BaseTool, error)
}

// 工具构造函数
type ToolConstructFunc func(ctx context.Context, toolConfig *ToolConfig, specificConfig interface{}) ([]tool.BaseTool, error)

// 工具构造器注册中心（类型+名称唯一，大小写无感）
var toolConstructorRegistry = NewComponentRegistry[*ToolDescriber, ToolConstructor]()

// 获取工具构造器
func GetToolConstructor(toolDesc *ToolDescriber) (ToolConstructor, error) {
	return toolConstructorRegistry.GetWithFallback(
		toolDesc,
		ToolNameGeneral,
		func(d *ToolDescriber) string { return d.Name },
		func(d *ToolDescriber) *ToolDescriber {
			return NewToolDescriber(ToolType(d.Type), ToolNameGeneral)
		},
	)
}

// 注册工具构造器
func RegisterToolConstructor(toolDesc *ToolDescriber, toolConstructor ToolConstructor, bindValues ...interface{}) error {
	return toolConstructorRegistry.Register(toolDesc, toolConstructor, bindValues...)
}

// 注册工具构造函数
func RegisterToolConstructFunc(toolType ToolType, toolName string, toolConstructFunc ToolConstructFunc, bindValues ...interface{}) error {
	toolDesc := NewToolDescriber(toolType, toolName)
	toolConstructor := NewComponentConstructor[*ToolDescriber, *ToolConfig, []tool.BaseTool](toolDesc, toolConstructFunc)
	return RegisterToolConstructor(toolDesc, toolConstructor, bindValues...)
}

func LookupToolDescriber(value interface{}) (*ToolDescriber, error) {
	return toolConstructorRegistry.LookupDesc(value)
}

func GetTool(ctx context.Context, toolOptions ...ToolOption) ([]tool.BaseTool, []*schema.ToolInfo, error) {
	toolConfig := NewToolConfig(toolOptions...)
	toolConstructor, err := GetToolConstructor(NewToolDescriber(ToolType(toolConfig.Type), toolConfig.Name))
	if err != nil {
		return nil, nil, err
	}
	tools, err := toolConstructor.Construct(ctx, toolConfig)
	if err != nil {
		return nil, nil, err
	}
	outTools := make([]tool.BaseTool, 0, len(tools))
	infos := make([]*schema.ToolInfo, 0, len(tools))
	for _, t := range tools {
		if t == nil {
			continue
		}
		outTools = append(outTools, t)
		info, infoErr := t.Info(ctx)
		if infoErr != nil {
			logger.Warnf("get tool info failed: %v", infoErr)
		}
		infos = append(infos, info) // nil表示没有工具信息
	}
	return outTools, infos, nil
}
