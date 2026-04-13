// 实现动态工具选择
package toolsearch

import (
	"context"

	"github.com/cloudwego/eino/adk"
	tsmiddleware "github.com/cloudwego/eino/adk/middlewares/dynamictool/toolsearch"
	"github.com/cloudwego/eino/components/tool"
	"github.com/leitiannet/einolib"
)

const (
	MiddlewareTypeToolSearch einolib.MiddlewareType = "toolsearch"
)

type ToolSearchMiddlewareConfig struct {
	tsmiddleware.Config // 内嵌结构体
}

func NewToolSearchMiddlewareConfig(toolSearchMiddlewareOptions ...ToolSearchMiddlewareOption) *ToolSearchMiddlewareConfig {
	config := &ToolSearchMiddlewareConfig{}
	einolib.ApplyOptions(config, toolSearchMiddlewareOptions)
	return config
}

type ToolSearchMiddlewareOption func(*ToolSearchMiddlewareConfig)

var (
	WithDynamicTools       = einolib.MakeAppendOption(func(c *ToolSearchMiddlewareConfig) *[]tool.BaseTool { return &c.DynamicTools })
	WithUseModelToolSearch = einolib.MakeOption(func(c *ToolSearchMiddlewareConfig, v bool) { c.UseModelToolSearch = v })
)

func NewToolSearchMiddleware(ctx context.Context, config *ToolSearchMiddlewareConfig) (adk.ChatModelAgentMiddleware, error) {
	return tsmiddleware.New(ctx, &config.Config)
}

func createMiddleware(ctx context.Context, middlewareConfig *einolib.MiddlewareConfig, specificConfig interface{}) (adk.ChatModelAgentMiddleware, error) {
	toolSearchMiddlewareConfig, err := einolib.ParseSpecificConfig(specificConfig, func() *ToolSearchMiddlewareConfig { return NewToolSearchMiddlewareConfig() })
	if err != nil {
		return nil, err
	}
	return NewToolSearchMiddleware(ctx, toolSearchMiddlewareConfig)
}

func init() {
	if err := einolib.RegisterMiddlewareConstructFunc(MiddlewareTypeToolSearch, einolib.GeneralMiddlewareName, createMiddleware, (*ToolSearchMiddlewareConfig)(nil)); err != nil {
		einolib.GetLogger().Errorf("register middleware %s failed: %v", MiddlewareTypeToolSearch, err)
	}
}
