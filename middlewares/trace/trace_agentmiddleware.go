package trace

import (
	"context"

	"github.com/cloudwego/eino/adk"
	"github.com/leitiannet/einolib"
)

func NewAgentMiddleware() adk.AgentMiddleware {
	return adk.AgentMiddleware{
		BeforeChatModel: func(ctx context.Context, state *adk.ChatModelAgentState) error {
			_, _ = ctx, state
			einolib.GetLogger().Infof("[AgentMiddleware] BeforeChatModel")
			return nil
		},
		AfterChatModel: func(ctx context.Context, state *adk.ChatModelAgentState) error {
			_, _ = ctx, state
			einolib.GetLogger().Infof("[AgentMiddleware] AfterChatModel")
			return nil
		},
	}
}
