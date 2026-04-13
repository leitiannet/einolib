package memory

import (
	"context"
	"fmt"

	"github.com/cloudwego/eino/adk/filesystem"
	"github.com/cloudwego/eino/schema"
)

type MemoryOperator struct {
	*filesystem.InMemoryBackend // 继承所有文件操作方法
}

func (o *MemoryOperator) Execute(ctx context.Context, input *filesystem.ExecuteRequest) (*filesystem.ExecuteResponse, error) {
	return nil, fmt.Errorf("memory operator does not support command execution: %s", input.Command)
}

func (o *MemoryOperator) ExecuteStreaming(ctx context.Context, input *filesystem.ExecuteRequest) (*schema.StreamReader[*filesystem.ExecuteResponse], error) {
	return nil, fmt.Errorf("memory operator does not support streaming command execution: %s", input.Command)
}
