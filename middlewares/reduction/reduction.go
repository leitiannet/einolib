// 控制工具结果占用的token数量（截断/清理）
package reduction

import (
	"context"

	"github.com/cloudwego/eino/adk"
	rdmiddleware "github.com/cloudwego/eino/adk/middlewares/reduction"
	"github.com/cloudwego/eino/schema"
	"github.com/leitiannet/einolib"
)

const (
	MiddlewareTypeReduction einolib.MiddlewareType = "reduction"
)

type ReductionMiddlewareConfig struct {
	rdmiddleware.Config // 内嵌结构体
}

func NewReductionMiddlewareConfig(reductionMiddlewareOptions ...ReductionMiddlewareOption) *ReductionMiddlewareConfig {
	reductionMiddlewareConfig := &ReductionMiddlewareConfig{}
	einolib.ApplyOptions(reductionMiddlewareConfig, reductionMiddlewareOptions)
	return reductionMiddlewareConfig
}

type ReductionMiddlewareOption func(*ReductionMiddlewareConfig)

var (
	WithBackend           = einolib.MakeOption(func(c *ReductionMiddlewareConfig, v rdmiddleware.Backend) { c.Backend = v })
	WithSkipTruncation    = einolib.MakeOption(func(c *ReductionMiddlewareConfig, v bool) { c.SkipTruncation = v })
	WithSkipClear         = einolib.MakeOption(func(c *ReductionMiddlewareConfig, v bool) { c.SkipClear = v })
	WithReadFileToolName  = einolib.MakeOption(func(c *ReductionMiddlewareConfig, v string) { c.ReadFileToolName = v })
	WithRootDir           = einolib.MakeOption(func(c *ReductionMiddlewareConfig, v string) { c.RootDir = v })
	WithMaxLengthForTrunc = einolib.MakeOption(func(c *ReductionMiddlewareConfig, v int) { c.MaxLengthForTrunc = v })
	WithTokenCounter      = einolib.MakeOption(func(c *ReductionMiddlewareConfig, v func(context.Context, []adk.Message, []*schema.ToolInfo) (int64, error)) {
		c.TokenCounter = v
	})
	WithMaxTokensForClear    = einolib.MakeOption(func(c *ReductionMiddlewareConfig, v int64) { c.MaxTokensForClear = v })
	WithClearRetentionSuffix = einolib.MakeOption(func(c *ReductionMiddlewareConfig, v int) { c.ClearRetentionSuffixLimit = v })
	WithClearAtLeastTokens   = einolib.MakeOption(func(c *ReductionMiddlewareConfig, v int64) { c.ClearAtLeastTokens = v })
	WithClearExcludeTools    = einolib.MakeAppendOption(func(c *ReductionMiddlewareConfig) *[]string { return &c.ClearExcludeTools })
	WithClearPostProcess     = einolib.MakeOption(func(c *ReductionMiddlewareConfig, v func(context.Context, *adk.ChatModelAgentState) context.Context) {
		c.ClearPostProcess = v
	})
	WithToolConfig = einolib.MakeOption(func(c *ReductionMiddlewareConfig, v map[string]*rdmiddleware.ToolReductionConfig) {
		c.ToolConfig = v
	})
)

func NewReductionMiddleware(ctx context.Context, reductionMiddlewareConfig *ReductionMiddlewareConfig) (adk.ChatModelAgentMiddleware, error) {
	return rdmiddleware.New(ctx, &reductionMiddlewareConfig.Config)
}

func createReductionMiddleware(ctx context.Context, middlewareConfig *einolib.MiddlewareConfig, specificConfig interface{}) (adk.ChatModelAgentMiddleware, error) {
	reductionMiddlewareConfig, err := einolib.ParseSpecificConfig(specificConfig, func() *ReductionMiddlewareConfig { return NewReductionMiddlewareConfig() })
	if err != nil {
		return nil, err
	}
	return NewReductionMiddleware(ctx, reductionMiddlewareConfig)
}

func init() {
	if err := einolib.RegisterMiddlewareConstructFunc(MiddlewareTypeReduction, einolib.GeneralMiddlewareName, createReductionMiddleware, (*ReductionMiddlewareConfig)(nil)); err != nil {
		einolib.GetLogger().Errorf("register middleware %s failed: %v", MiddlewareTypeReduction, err)
	}
}
