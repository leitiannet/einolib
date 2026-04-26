// 自动压缩对话历史
package summarization

import (
	"context"

	"github.com/cloudwego/eino/adk"
	summiddleware "github.com/cloudwego/eino/adk/middlewares/summarization"
	"github.com/cloudwego/eino/components/model"
	"github.com/leitiannet/einolib"
)

type SummarizationMiddlewareConfig struct {
	summiddleware.Config // 内嵌结构体
}

func NewSummarizationMiddlewareConfig(summarizationMiddlewareOptions ...SummarizationMiddlewareOption) *SummarizationMiddlewareConfig {
	summarizationMiddlewareConfig := &SummarizationMiddlewareConfig{}
	einolib.ApplyOptions(summarizationMiddlewareConfig, summarizationMiddlewareOptions)
	return summarizationMiddlewareConfig
}

type SummarizationMiddlewareOption func(*SummarizationMiddlewareConfig)

var (
	WithModel                = einolib.MakeOption(func(c *SummarizationMiddlewareConfig, v model.BaseChatModel) { c.Model = v })
	WithModelOptions         = einolib.MakeAppendOption(func(c *SummarizationMiddlewareConfig) *[]model.Option { return &c.ModelOptions })
	WithTokenCounter         = einolib.MakeOption(func(c *SummarizationMiddlewareConfig, v summiddleware.TokenCounterFunc) { c.TokenCounter = v })
	WithTrigger              = einolib.MakeOption(func(c *SummarizationMiddlewareConfig, v *summiddleware.TriggerCondition) { c.Trigger = v })
	WithEmitInternalEvents   = einolib.MakeOption(func(c *SummarizationMiddlewareConfig, v bool) { c.EmitInternalEvents = v })
	WithUserInstruction      = einolib.MakeOption(func(c *SummarizationMiddlewareConfig, v string) { c.UserInstruction = v })
	WithTranscriptFilePath   = einolib.MakeOption(func(c *SummarizationMiddlewareConfig, v string) { c.TranscriptFilePath = v })
	WithGenModelInput        = einolib.MakeOption(func(c *SummarizationMiddlewareConfig, v summiddleware.GenModelInputFunc) { c.GenModelInput = v })
	WithFinalize             = einolib.MakeOption(func(c *SummarizationMiddlewareConfig, v summiddleware.FinalizeFunc) { c.Finalize = v })
	WithCallback             = einolib.MakeOption(func(c *SummarizationMiddlewareConfig, v summiddleware.CallbackFunc) { c.Callback = v })
	WithPreserveUserMessages = einolib.MakeOption(func(c *SummarizationMiddlewareConfig, v *summiddleware.PreserveUserMessages) {
		c.PreserveUserMessages = v
	})
)

func NewSummarizationMiddleware(ctx context.Context, summarizationMiddlewareConfig *SummarizationMiddlewareConfig) (adk.ChatModelAgentMiddleware, error) {
	return summiddleware.New(ctx, &summarizationMiddlewareConfig.Config)
}

func createSummarizationMiddleware(ctx context.Context, middlewareConfig *einolib.MiddlewareConfig, specificConfig interface{}) (adk.ChatModelAgentMiddleware, error) {
	summarizationMiddlewareConfig, err := einolib.ParseSpecificConfig(specificConfig, func() *SummarizationMiddlewareConfig { return NewSummarizationMiddlewareConfig() })
	if err != nil {
		return nil, err
	}
	return NewSummarizationMiddleware(ctx, summarizationMiddlewareConfig)
}

func init() {
	if err := einolib.RegisterMiddlewareConstructFunc(einolib.MiddlewareTypeSummarization, einolib.MiddlewareNameGeneral, createSummarizationMiddleware, (*SummarizationMiddlewareConfig)(nil)); err != nil {
		einolib.GetLogger().Errorf("register middleware %s failed: %v", einolib.MiddlewareTypeSummarization, err)
	}
}
