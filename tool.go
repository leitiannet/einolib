package einolib

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/schema"
)

// 工具类型
type ToolType string

const (
	ToolTypeUnknown ToolType = ""        // 未知类型工具
	ToolTypeCustom  ToolType = "custom"  // 用户定义工具
	ToolTypeBuiltin ToolType = "builtin" // 官方内置工具
	ToolTypeMCP     ToolType = "mcp"     // MCP协议工具
)

// 合法的工具类型列表
var validToolTypes = []ToolType{ToolTypeCustom, ToolTypeBuiltin, ToolTypeMCP}

var validToolTypeMap = func() map[ToolType]bool {
	m := make(map[ToolType]bool, len(validToolTypes)+1)
	for _, t := range validToolTypes {
		m[t] = true
	}
	m[ToolTypeUnknown] = false
	return m
}()

const GeneralToolName = "*" // 通用工具名称

// 工具描述
type ToolDescriber struct {
	ToolType ToolType // 工具类型
	ToolName string   // 工具名称
}

func NewToolDescriber(toolType ToolType, toolName string) *ToolDescriber {
	return &ToolDescriber{
		ToolType: toolType,
		ToolName: toolName,
	}
}

func (td *ToolDescriber) String() string {
	return fmt.Sprintf("%s:%s", td.ToolType, td.ToolName)
}

func (td *ToolDescriber) Key() string {
	return strings.ToLower(td.String())
}

func (td *ToolDescriber) Validate() error {
	if valid, ok := validToolTypeMap[td.ToolType]; !ok || !valid {
		return fmt.Errorf("toolType invalid: %q", td.ToolType)
	}
	if td.ToolName == "" {
		return fmt.Errorf("toolName invalid: %q", td.ToolName)
	}
	return nil
}

// 工具配置
type ToolConfig struct {
	ComponentConfig          // 特定配置
	schema.ToolInfo          // 工具信息
	ToolType        ToolType // 工具类型
}

func NewToolConfig(toolOptions ...ToolOption) *ToolConfig {
	toolConfig := &ToolConfig{
		ComponentConfig: ComponentConfig{
			ConfigMap: NewSyncMap(),
		},
	}
	ApplyOptions(toolConfig, toolOptions)
	return toolConfig
}

// 工具选项
type ToolOption func(toolConfig *ToolConfig)

var (
	WithToolType            = MakeOption(func(c *ToolConfig, v ToolType) { c.ToolType = v })
	WithToolName            = MakeOption(func(c *ToolConfig, v string) { c.Name = v })
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
	constructor, err := toolConstructorRegistry.Get(toolDesc)
	if err != nil && toolDesc.ToolName != GeneralToolName {
		generalConstructor, generalErr := toolConstructorRegistry.Get(NewToolDescriber(toolDesc.ToolType, GeneralToolName))
		if generalErr != nil {
			return constructor, err
		}
		return generalConstructor, nil
	}
	return constructor, err
}

// 注册工具构造器
func RegisterToolConstructor(toolDesc *ToolDescriber, toolConstructor ToolConstructor) error {
	return toolConstructorRegistry.Register(toolDesc, toolConstructor)
}

// 注册工具构造函数；无特定配置时不要传参
func RegisterToolConstructFunc(toolType ToolType, toolName string, toolConstructFunc ToolConstructFunc, bindValues ...interface{}) error {
	toolDesc := NewToolDescriber(toolType, toolName)
	toolConstructor := NewComponentConstructor[*ToolDescriber, *ToolConfig, []tool.BaseTool](toolDesc, toolConstructFunc)
	return toolConstructorRegistry.Register(toolDesc, toolConstructor, bindValues...)
}

func LookupToolDescriber(value interface{}) (*ToolDescriber, error) {
	return toolConstructorRegistry.LookupDesc(value)
}

func GetTool(ctx context.Context, toolOptions ...ToolOption) ([]tool.BaseTool, []*schema.ToolInfo, error) {
	toolConfig := NewToolConfig(toolOptions...)
	if toolConfig.ToolType != ToolTypeUnknown {
		return getToolByType(ctx, toolConfig.ToolType, toolConfig)
	}
	// 并发获取所有类型工具
	var (
		wg        sync.WaitGroup
		allToolMu sync.Mutex
		allTools  []tool.BaseTool
		allInfos  []*schema.ToolInfo
		allErrs   []error
	)
	for _, toolType := range validToolTypes {
		wg.Add(1)
		go func(t ToolType) {
			defer wg.Done()
			typeTools, typeInfos, err := getToolByType(ctx, t, toolConfig)
			allToolMu.Lock()
			defer allToolMu.Unlock()
			if err != nil {
				allErrs = append(allErrs, fmt.Errorf("[%s] %v", t, err))
				return
			}
			if len(typeTools) > 0 {
				allTools = append(allTools, typeTools...)
				allInfos = append(allInfos, typeInfos...)
			}
		}(toolType)
	}
	wg.Wait()
	// 只要拿到一个工具即视为成功（部分类型失败静默忽略）
	if len(allTools) > 0 {
		return allTools, allInfos, nil
	}
	if len(allErrs) > 0 {
		return nil, nil, errors.Join(allErrs...)
	}
	return nil, nil, nil
}

func getToolByType(ctx context.Context, toolType ToolType, toolConfig *ToolConfig) ([]tool.BaseTool, []*schema.ToolInfo, error) {
	toolConstructor, err := GetToolConstructor(NewToolDescriber(toolType, toolConfig.Name))
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
