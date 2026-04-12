// 观察执行顺序（日志输出，不修改上下文与返回值）
package trace

import (
	"context"

	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/components/model"
	"github.com/leitiannet/einolib"
)

// 返回结构体形式的AgentMiddleware
func AgentMiddleware() adk.AgentMiddleware {
	return adk.AgentMiddleware{
		BeforeChatModel: func(_ context.Context, _ *adk.ChatModelAgentState) error {
			einolib.GetLogger().Infof("[AgentMiddleware] BeforeChatModel")
			return nil
		},
		AfterChatModel: func(_ context.Context, _ *adk.ChatModelAgentState) error {
			einolib.GetLogger().Infof("[AgentMiddleware] AfterChatModel")
			return nil
		},
	}
}

// 实现adk.ChatModelAgentMiddleware，覆盖全部钩子，仅打印日志后透传
type ChatModelAgentMiddleware struct{}

func NewChatModelAgentMiddleware() adk.ChatModelAgentMiddleware {
	return &ChatModelAgentMiddleware{}
}

func (*ChatModelAgentMiddleware) BeforeAgent(ctx context.Context, runCtx *adk.ChatModelAgentContext) (context.Context, *adk.ChatModelAgentContext, error) {
	einolib.GetLogger().Infof("[ChatModelAgentMiddleware] BeforeAgent")
	return ctx, runCtx, nil
}

func (*ChatModelAgentMiddleware) BeforeModelRewriteState(ctx context.Context, state *adk.ChatModelAgentState, mc *adk.ModelContext) (context.Context, *adk.ChatModelAgentState, error) {
	einolib.GetLogger().Infof("[ChatModelAgentMiddleware] BeforeModelRewriteState")
	return ctx, state, nil
}

func (*ChatModelAgentMiddleware) AfterModelRewriteState(ctx context.Context, state *adk.ChatModelAgentState, mc *adk.ModelContext) (context.Context, *adk.ChatModelAgentState, error) {
	einolib.GetLogger().Infof("[ChatModelAgentMiddleware] AfterModelRewriteState")
	return ctx, state, nil
}

func (*ChatModelAgentMiddleware) WrapModel(ctx context.Context, m model.BaseChatModel, mc *adk.ModelContext) (model.BaseChatModel, error) {
	einolib.GetLogger().Infof("[ChatModelAgentMiddleware] WrapModel")
	return m, nil
}

func (*ChatModelAgentMiddleware) WrapInvokableToolCall(ctx context.Context, endpoint adk.InvokableToolCallEndpoint, tCtx *adk.ToolContext) (adk.InvokableToolCallEndpoint, error) {
	einolib.GetLogger().Infof("[ChatModelAgentMiddleware] WrapInvokableToolCall tool=%q", toolName(tCtx))
	return endpoint, nil
}

func (*ChatModelAgentMiddleware) WrapStreamableToolCall(ctx context.Context, endpoint adk.StreamableToolCallEndpoint, tCtx *adk.ToolContext) (adk.StreamableToolCallEndpoint, error) {
	einolib.GetLogger().Infof("[ChatModelAgentMiddleware] WrapStreamableToolCall tool=%q", toolName(tCtx))
	return endpoint, nil
}

func (*ChatModelAgentMiddleware) WrapEnhancedInvokableToolCall(ctx context.Context, endpoint adk.EnhancedInvokableToolCallEndpoint, tCtx *adk.ToolContext) (adk.EnhancedInvokableToolCallEndpoint, error) {
	einolib.GetLogger().Infof("[ChatModelAgentMiddleware] WrapEnhancedInvokableToolCall tool=%q", toolName(tCtx))
	return endpoint, nil
}

func (*ChatModelAgentMiddleware) WrapEnhancedStreamableToolCall(ctx context.Context, endpoint adk.EnhancedStreamableToolCallEndpoint, tCtx *adk.ToolContext) (adk.EnhancedStreamableToolCallEndpoint, error) {
	einolib.GetLogger().Infof("[ChatModelAgentMiddleware] WrapEnhancedStreamableToolCall tool=%q", toolName(tCtx))
	return endpoint, nil
}

func toolName(tCtx *adk.ToolContext) string {
	if tCtx == nil {
		return ""
	}
	return tCtx.Name
}
