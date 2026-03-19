package einolib

import (
	"context"
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
	ToolTypeAll     ToolType = "*"       // 所有类型工具
)

// 合法的工具类型列表
var validToolTypes = []ToolType{ToolTypeCustom, ToolTypeBuiltin, ToolTypeMCP}

var validToolTypeMap = func() map[ToolType]bool {
	m := make(map[ToolType]bool, len(validToolTypes)+2)
	for _, t := range validToolTypes {
		m[t] = true
	}
	m[ToolTypeUnknown] = false
	m[ToolTypeAll] = false
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
		return fmt.Errorf("toolType invalid: %s", td.ToolType)
	}
	return nil
}

// 工具配置
type ToolConfig struct {
	ComponentConfig          // 特定配置
	schema.ToolInfo          // 工具信息
	ToolType        ToolType // 工具类型
}

func NewToolConfig(toolType ToolType, toolName string, toolOptions ...ToolOption) *ToolConfig {
	toolConfig := &ToolConfig{
		ComponentConfig: ComponentConfig{
			ConfigMap: NewSyncMap(),
		},
		ToolInfo: schema.ToolInfo{
			Name: toolName,
		},
		ToolType: toolType,
	}
	ApplyOptions(toolConfig, toolOptions)
	return toolConfig
}

// 工具选项
type ToolOption func(toolConfig *ToolConfig)

func WithToolComponentConfig(toolType ToolType, toolName string, value interface{}) ToolOption {
	return func(toolConfig *ToolConfig) {
		if toolConfig != nil {
			toolConfig.SetConfig(NewToolDescriber(toolType, toolName), value)
		}
	}
}

func WithToolType(toolType ToolType) ToolOption {
	return func(toolConfig *ToolConfig) {
		if toolConfig != nil {
			toolConfig.ToolType = toolType
		}
	}
}

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
	return toolConstructorRegistry.Get(toolDesc)
}

// 注册工具构造器
func RegisterToolConstructor(toolDesc *ToolDescriber, toolConstructor ToolConstructor) error {
	return toolConstructorRegistry.Register(toolDesc, toolConstructor)
}

// 注册工具构造函数
func RegisterToolConstructFunc(toolType ToolType, toolName string, toolConstructFunc ToolConstructFunc) error {
	toolDesc := NewToolDescriber(toolType, toolName)
	toolConstructor := NewComponentConstructor[*ToolDescriber, *ToolConfig, []tool.BaseTool](toolDesc, toolConstructFunc)
	return RegisterToolConstructor(toolDesc, toolConstructor)
}

func GetTool(ctx context.Context, toolType ToolType, toolName string, toolOptions ...ToolOption) ([]tool.BaseTool, []*schema.ToolInfo, error) {
	toolConfig := NewToolConfig(toolType, toolName, toolOptions...)
	if toolConfig.ToolType != ToolTypeUnknown {
		return getToolByType(ctx, toolConfig.ToolType, toolConfig)
	}
	// 并发获取所有类型工具
	var (
		wg        sync.WaitGroup
		allToolMu sync.Mutex
		allTools  []tool.BaseTool
		allInfos  []*schema.ToolInfo
		firstErr  error
	)
	for _, toolType := range validToolTypes {
		wg.Add(1)
		go func(t ToolType) {
			defer wg.Done()
			typeTools, typeInfos, err := getToolByType(ctx, t, toolConfig)
			allToolMu.Lock()
			defer allToolMu.Unlock()
			if err != nil {
				if firstErr == nil {
					firstErr = err
				}
				return
			}
			if len(typeTools) > 0 {
				allTools = append(allTools, typeTools...)
				allInfos = append(allInfos, typeInfos...)
			}
		}(toolType)
	}
	wg.Wait()
	if len(allTools) == 0 && firstErr != nil {
		return nil, nil, firstErr
	}
	return allTools, allInfos, nil
}

func getToolByType(ctx context.Context, toolType ToolType, toolConfig *ToolConfig) ([]tool.BaseTool, []*schema.ToolInfo, error) {
	toolConstructor, err := GetToolConstructor(NewToolDescriber(toolType, toolConfig.Name))
	if err != nil {
		if toolType == ToolTypeMCP && toolConfig.Name != "" {
			toolConstructor, err = GetToolConstructor(NewToolDescriber(toolType, GeneralToolName))
		}
		if err != nil {
			return nil, nil, err
		}
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
		info, _ := t.Info(ctx) // 可能为nil
		infos = append(infos, info)
	}
	return outTools, infos, nil
}
