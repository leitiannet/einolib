// 修复消息历史中悬空的工具调用，为缺少响应的工具调用自动插入占位符消息
package patchtoolcalls

import (
	"context"

	"github.com/cloudwego/eino/adk"
	ptcmiddleware "github.com/cloudwego/eino/adk/middlewares/patchtoolcalls"
	"github.com/leitiannet/einolib"
)

const (
	MiddlewareTypePatchToolCalls einolib.MiddlewareType = "patchtoolcalls"
)

type PatchToolCallsMiddlewareConfig struct {
	ptcmiddleware.Config // 内嵌结构体
}

func NewPatchToolCallsMiddlewareConfig(patchToolCallsMiddlewareOptions ...PatchToolCallsMiddlewareOption) *PatchToolCallsMiddlewareConfig {
	config := &PatchToolCallsMiddlewareConfig{}
	einolib.ApplyOptions(config, patchToolCallsMiddlewareOptions)
	return config
}

type PatchToolCallsMiddlewareOption func(*PatchToolCallsMiddlewareConfig)

var (
	WithPatchedContentGenerator = einolib.MakeOption(func(c *PatchToolCallsMiddlewareConfig, v func(ctx context.Context, toolName, toolCallID string) (string, error)) {
		c.PatchedContentGenerator = v
	})
)

func NewPatchToolCallsMiddleware(ctx context.Context, config *PatchToolCallsMiddlewareConfig) (adk.ChatModelAgentMiddleware, error) {
	return ptcmiddleware.New(ctx, &config.Config)
}

func createMiddleware(ctx context.Context, middlewareConfig *einolib.MiddlewareConfig, specificConfig interface{}) (adk.ChatModelAgentMiddleware, error) {
	patchToolCallsMiddlewareConfig, err := einolib.ParseSpecificConfig(specificConfig, func() *PatchToolCallsMiddlewareConfig { return NewPatchToolCallsMiddlewareConfig() })
	if err != nil {
		return nil, err
	}
	return NewPatchToolCallsMiddleware(ctx, patchToolCallsMiddlewareConfig)
}

func init() {
	if err := einolib.RegisterMiddlewareConstructFunc(MiddlewareTypePatchToolCalls, einolib.GeneralMiddlewareName, createMiddleware, (*PatchToolCallsMiddlewareConfig)(nil)); err != nil {
		einolib.GetLogger().Errorf("register middleware %s failed: %v", MiddlewareTypePatchToolCalls, err)
	}
}
