// 实现动态工具选择
package toolsearch

import (
	"context"

	"github.com/cloudwego/eino/adk"
	tsmiddleware "github.com/cloudwego/eino/adk/middlewares/dynamictool/toolsearch"
	"github.com/cloudwego/eino/components/tool"
	"github.com/leitiannet/einolib"
)

type ToolSearchMiddlewareConfig struct {
	tsmiddleware.Config // 内嵌结构体
}

func NewToolSearchMiddlewareConfig(toolSearchMiddlewareOptions ...ToolSearchMiddlewareOption) *ToolSearchMiddlewareConfig {
	toolSearchMiddlewareConfig := &ToolSearchMiddlewareConfig{}
	einolib.ApplyOptions(toolSearchMiddlewareConfig, toolSearchMiddlewareOptions)
	return toolSearchMiddlewareConfig
}

type ToolSearchMiddlewareOption func(*ToolSearchMiddlewareConfig)

var (
	WithDynamicTools       = einolib.MakeAppendOption(func(c *ToolSearchMiddlewareConfig) *[]tool.BaseTool { return &c.DynamicTools })
	WithUseModelToolSearch = einolib.MakeOption(func(c *ToolSearchMiddlewareConfig, v bool) { c.UseModelToolSearch = v })
)

func NewToolSearchMiddleware(ctx context.Context, toolSearchMiddlewareConfig *ToolSearchMiddlewareConfig) (adk.ChatModelAgentMiddleware, error) {
	return tsmiddleware.New(ctx, &toolSearchMiddlewareConfig.Config)
}

func createToolSearchMiddleware(ctx context.Context, middlewareConfig *einolib.MiddlewareConfig, specificConfig interface{}) (adk.ChatModelAgentMiddleware, error) {
	toolSearchMiddlewareConfig, err := einolib.ParseSpecificConfig(specificConfig, func() *ToolSearchMiddlewareConfig { return NewToolSearchMiddlewareConfig() })
	if err != nil {
		return nil, err
	}
	return NewToolSearchMiddleware(ctx, toolSearchMiddlewareConfig)
}

func init() {
	if err := einolib.RegisterMiddlewareConstructFunc(einolib.MiddlewareTypeToolSearch, einolib.MiddlewareNameGeneral, createToolSearchMiddleware, (*ToolSearchMiddlewareConfig)(nil)); err != nil {
		einolib.GetLogger().Errorf("register middleware %s failed: %v", einolib.MiddlewareTypeToolSearch, err)
	}
}
