package agentkit

import (
	"context"
	"fmt"

	einoagentkit "github.com/cloudwego/eino-ext/adk/backend/agentkit"
	"github.com/cloudwego/eino/adk/filesystem"
	"github.com/cloudwego/eino/schema"
)

type AgentKitOperator struct {
	*einoagentkit.SandboxTool // 继承所有文件操作方法
}

func (o *AgentKitOperator) ExecuteStreaming(_ context.Context, input *filesystem.ExecuteRequest) (*schema.StreamReader[*filesystem.ExecuteResponse], error) {
	return nil, fmt.Errorf("agentkit operator does not support streaming command execution: %s", input.Command)
}
