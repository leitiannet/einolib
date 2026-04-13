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

func (m *ChatModelAgentMiddleware) BeforeModelRewriteState(ctx context.Context, state *adk.ChatModelAgentState, _ *adk.ModelContext) (context.Context, *adk.ChatModelAgentState, error) {
	einolib.GetLogger().Infof("[%s] BeforeModelRewriteState", m.prefix)
	return ctx, state, nil
}

func (m *ChatModelAgentMiddleware) AfterModelRewriteState(ctx context.Context, state *adk.ChatModelAgentState, _ *adk.ModelContext) (context.Context, *adk.ChatModelAgentState, error) {
	einolib.GetLogger().Infof("[%s] AfterModelRewriteState", m.prefix)
	return ctx, state, nil
}

func (m *ChatModelAgentMiddleware) AfterToolCallsRewriteState(ctx context.Context, state *adk.ChatModelAgentState, tc *adk.ToolCallsContext) (context.Context, *adk.ChatModelAgentState, error) {
	einolib.GetLogger().Infof("[%s] AfterToolCallsRewriteState", m.prefix)
	return ctx, state, nil
}

func (m *ChatModelAgentMiddleware) WrapModel(_ context.Context, cm model.BaseChatModel, _ *adk.ModelContext) (model.BaseChatModel, error) {
	einolib.GetLogger().Infof("[%s] WrapModel", m.prefix)
	return cm, nil
}

func (m *ChatModelAgentMiddleware) WrapInvokableToolCall(_ context.Context, endpoint adk.InvokableToolCallEndpoint, tCtx *adk.ToolContext) (adk.InvokableToolCallEndpoint, error) {
	einolib.GetLogger().Infof("[%s] WrapInvokableToolCall tool=%q", m.prefix, toolName(tCtx))
	return endpoint, nil
}

func (m *ChatModelAgentMiddleware) WrapStreamableToolCall(_ context.Context, endpoint adk.StreamableToolCallEndpoint, tCtx *adk.ToolContext) (adk.StreamableToolCallEndpoint, error) {
	einolib.GetLogger().Infof("[%s] WrapStreamableToolCall tool=%q", m.prefix, toolName(tCtx))
	return endpoint, nil
}

func (m *ChatModelAgentMiddleware) WrapEnhancedInvokableToolCall(_ context.Context, endpoint adk.EnhancedInvokableToolCallEndpoint, tCtx *adk.ToolContext) (adk.EnhancedInvokableToolCallEndpoint, error) {
	einolib.GetLogger().Infof("[%s] WrapEnhancedInvokableToolCall tool=%q", m.prefix, toolName(tCtx))
	return endpoint, nil
}

func (m *ChatModelAgentMiddleware) WrapEnhancedStreamableToolCall(_ context.Context, endpoint adk.EnhancedStreamableToolCallEndpoint, tCtx *adk.ToolContext) (adk.EnhancedStreamableToolCallEndpoint, error) {
	einolib.GetLogger().Infof("[%s] WrapEnhancedStreamableToolCall tool=%q", m.prefix, toolName(tCtx))
	return endpoint, nil
}

func toolName(tCtx *adk.ToolContext) string {
	if tCtx == nil {
		return ""
	}
	return tCtx.Name
}
