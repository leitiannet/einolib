// 观察执行顺序（日志输出，不修改上下文与返回值）
package trace

import (
	"context"

	"github.com/cloudwego/eino/adk"
	"github.com/leitiannet/einolib"
)

const (
	MiddlewareTypeTrace einolib.MiddlewareType = "trace"
)

type TraceMiddlewareConfig struct {
	Prefix string
}

func NewTraceMiddlewareConfig(traceMiddlewareOptions ...TraceMiddlewareOption) *TraceMiddlewareConfig {
	traceMiddlewareConfig := &TraceMiddlewareConfig{Prefix: "trace"}
	einolib.ApplyOptions(traceMiddlewareConfig, traceMiddlewareOptions)
	return traceMiddlewareConfig
}

type TraceMiddlewareOption func(*TraceMiddlewareConfig)

var (
	WithPrefix = einolib.MakeOption(func(c *TraceMiddlewareConfig, v string) { c.Prefix = v })
)

func NewTraceMiddleware(ctx context.Context, traceMiddlewareConfig *TraceMiddlewareConfig) (*ChatModelAgentMiddleware, error) {
	_ = ctx
	return &ChatModelAgentMiddleware{prefix: traceMiddlewareConfig.Prefix}, nil
}

func createTraceMiddleware(ctx context.Context, middlewareConfig *einolib.MiddlewareConfig, specificConfig interface{}) (adk.ChatModelAgentMiddleware, error) {
	traceMiddlewareConfig, err := einolib.ParseSpecificConfig(specificConfig, func() *TraceMiddlewareConfig { return NewTraceMiddlewareConfig() })
	if err != nil {
		return nil, err
	}
	return NewTraceMiddleware(ctx, traceMiddlewareConfig)
}

func init() {
	if err := einolib.RegisterMiddlewareConstructFunc(MiddlewareTypeTrace, einolib.GeneralMiddlewareName, createTraceMiddleware, (*TraceMiddlewareConfig)(nil)); err != nil {
		einolib.GetLogger().Errorf("register middleware %s failed: %v", MiddlewareTypeTrace, err)
	}
}
