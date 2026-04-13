// 将工具调用错误转为字符串结果（或流式单条文本），并透传 interrupt-rerun 相关错误
package safetool

import (
	"context"

	"github.com/cloudwego/eino/adk"
	"github.com/leitiannet/einolib"
)

const (
	MiddlewareTypeSafeTool einolib.MiddlewareType = "safetool"
)

type SafeToolMiddlewareConfig struct {
	ErrorFormat string
}

func NewSafeToolMiddlewareConfig(safeToolMiddlewareOptions ...SafeToolMiddlewareOption) *SafeToolMiddlewareConfig {
	config := &SafeToolMiddlewareConfig{ErrorFormat: "[tool error] %v"}
	einolib.ApplyOptions(config, safeToolMiddlewareOptions)
	return config
}

type SafeToolMiddlewareOption func(*SafeToolMiddlewareConfig)

var (
	WithErrorFormat = einolib.MakeOption(func(c *SafeToolMiddlewareConfig, v string) { c.ErrorFormat = v })
)

func NewSafeToolMiddleware(_ context.Context, config *SafeToolMiddlewareConfig) (*ChatModelAgentMiddleware, error) {
	return &ChatModelAgentMiddleware{
		errorFormat: config.ErrorFormat,
	}, nil
}

func createMiddleware(ctx context.Context, middlewareConfig *einolib.MiddlewareConfig, specificConfig interface{}) (adk.ChatModelAgentMiddleware, error) {
	safeToolMiddlewareConfig, err := einolib.ParseSpecificConfig(specificConfig, func() *SafeToolMiddlewareConfig { return NewSafeToolMiddlewareConfig() })
	if err != nil {
		return nil, err
	}
	return NewSafeToolMiddleware(ctx, safeToolMiddlewareConfig)
}

func init() {
	if err := einolib.RegisterMiddlewareConstructFunc(MiddlewareTypeSafeTool, einolib.GeneralMiddlewareName, createMiddleware, (*SafeToolMiddlewareConfig)(nil)); err != nil {
		einolib.GetLogger().Errorf("register middleware %s failed: %v", MiddlewareTypeSafeTool, err)
	}
}
