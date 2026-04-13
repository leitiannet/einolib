package trace

import (
	"context"

	"github.com/cloudwego/eino/adk"
	"github.com/leitiannet/einolib"
)

func NewAgentMiddleware() adk.AgentMiddleware {
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
