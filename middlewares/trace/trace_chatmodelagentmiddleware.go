package trace

import (
	"context"

	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/components/model"
	"github.com/leitiannet/einolib"
)

type ChatModelAgentMiddleware struct {
	prefix string
}

func (m *ChatModelAgentMiddleware) BeforeAgent(ctx context.Context, runCtx *adk.ChatModelAgentContext) (context.Context, *adk.ChatModelAgentContext, error) {
	einolib.GetLogger().Infof("[%s] BeforeAgent", m.prefix)
	return ctx, runCtx, nil
}

func (m *ChatModelAgentMiddleware) BeforeModelRewriteState(ctx context.Context, state *adk.ChatModelAgentState, modelCtx *adk.ModelContext) (context.Context, *adk.ChatModelAgentState, error) {
	_ = modelCtx
	einolib.GetLogger().Infof("[%s] BeforeModelRewriteState", m.prefix)
	return ctx, state, nil
}

func (m *ChatModelAgentMiddleware) AfterModelRewriteState(ctx context.Context, state *adk.ChatModelAgentState, modelCtx *adk.ModelContext) (context.Context, *adk.ChatModelAgentState, error) {
	_ = modelCtx
	einolib.GetLogger().Infof("[%s] AfterModelRewriteState", m.prefix)
	return ctx, state, nil
}

func (m *ChatModelAgentMiddleware) AfterToolCallsRewriteState(ctx context.Context, state *adk.ChatModelAgentState, toolCallsCtx *adk.ToolCallsContext) (context.Context, *adk.ChatModelAgentState, error) {
	_ = toolCallsCtx
	einolib.GetLogger().Infof("[%s] AfterToolCallsRewriteState", m.prefix)
	return ctx, state, nil
}

func (m *ChatModelAgentMiddleware) WrapModel(ctx context.Context, cm model.BaseChatModel, modelCtx *adk.ModelContext) (model.BaseChatModel, error) {
	_, _ = ctx, modelCtx
	einolib.GetLogger().Infof("[%s] WrapModel", m.prefix)
	return cm, nil
}

func (m *ChatModelAgentMiddleware) WrapInvokableToolCall(ctx context.Context, endpoint adk.InvokableToolCallEndpoint, tCtx *adk.ToolContext) (adk.InvokableToolCallEndpoint, error) {
	_ = ctx
	einolib.GetLogger().Infof("[%s] WrapInvokableToolCall tool=%q", m.prefix, toolName(tCtx))
	return endpoint, nil
}

func (m *ChatModelAgentMiddleware) WrapStreamableToolCall(ctx context.Context, endpoint adk.StreamableToolCallEndpoint, tCtx *adk.ToolContext) (adk.StreamableToolCallEndpoint, error) {
	_ = ctx
	einolib.GetLogger().Infof("[%s] WrapStreamableToolCall tool=%q", m.prefix, toolName(tCtx))
	return endpoint, nil
}

func (m *ChatModelAgentMiddleware) WrapEnhancedInvokableToolCall(ctx context.Context, endpoint adk.EnhancedInvokableToolCallEndpoint, tCtx *adk.ToolContext) (adk.EnhancedInvokableToolCallEndpoint, error) {
	_ = ctx
	einolib.GetLogger().Infof("[%s] WrapEnhancedInvokableToolCall tool=%q", m.prefix, toolName(tCtx))
	return endpoint, nil
}

func (m *ChatModelAgentMiddleware) WrapEnhancedStreamableToolCall(ctx context.Context, endpoint adk.EnhancedStreamableToolCallEndpoint, tCtx *adk.ToolContext) (adk.EnhancedStreamableToolCallEndpoint, error) {
	_ = ctx
	einolib.GetLogger().Infof("[%s] WrapEnhancedStreamableToolCall tool=%q", m.prefix, toolName(tCtx))
	return endpoint, nil
}

func toolName(tCtx *adk.ToolContext) string {
	if tCtx == nil {
		return ""
	}
	return tCtx.Name
}
