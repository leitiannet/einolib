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
	config := &TraceMiddlewareConfig{Prefix: "trace"}
	einolib.ApplyOptions(config, traceMiddlewareOptions)
	return config
}

type TraceMiddlewareOption func(*TraceMiddlewareConfig)

var (
	WithPrefix = einolib.MakeOption(func(c *TraceMiddlewareConfig, v string) { c.Prefix = v })
)

func NewTraceMiddleware(_ context.Context, config *TraceMiddlewareConfig) (*ChatModelAgentMiddleware, error) {
	return &ChatModelAgentMiddleware{prefix: config.Prefix}, nil
}

func createMiddleware(ctx context.Context, middlewareConfig *einolib.MiddlewareConfig, specificConfig interface{}) (adk.ChatModelAgentMiddleware, error) {
	traceMiddlewareConfig, err := einolib.ParseSpecificConfig(specificConfig, func() *TraceMiddlewareConfig { return NewTraceMiddlewareConfig() })
	if err != nil {
		return nil, err
	}
	return NewTraceMiddleware(ctx, traceMiddlewareConfig)
}

func init() {
	if err := einolib.RegisterMiddlewareConstructFunc(MiddlewareTypeTrace, einolib.GeneralMiddlewareName, createMiddleware, (*TraceMiddlewareConfig)(nil)); err != nil {
		einolib.GetLogger().Errorf("register middleware %s failed: %v", MiddlewareTypeTrace, err)
	}
}
