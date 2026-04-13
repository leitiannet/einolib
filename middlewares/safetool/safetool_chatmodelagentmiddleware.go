package safetool

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
)

type ChatModelAgentMiddleware struct {
	*adk.BaseChatModelAgentMiddleware
	errorFormat string
}

func (m *ChatModelAgentMiddleware) WrapInvokableToolCall(_ context.Context, endpoint adk.InvokableToolCallEndpoint, _ *adk.ToolContext) (adk.InvokableToolCallEndpoint, error) {
	return func(ctx context.Context, argumentsInJSON string, opts ...tool.Option) (string, error) {
		result, err := endpoint(ctx, argumentsInJSON, opts...)
		if err != nil {
			if _, ok := compose.IsInterruptRerunError(err); ok {
				return result, err
			}
			return fmt.Sprintf(m.errorFormat, err), nil
		}
		return result, nil
	}, nil
}

func (m *ChatModelAgentMiddleware) WrapStreamableToolCall(_ context.Context, endpoint adk.StreamableToolCallEndpoint, _ *adk.ToolContext) (adk.StreamableToolCallEndpoint, error) {
	return func(ctx context.Context, argumentsInJSON string, opts ...tool.Option) (*schema.StreamReader[string], error) {
		sr, err := endpoint(ctx, argumentsInJSON, opts...)
		if err != nil {
			if _, ok := compose.IsInterruptRerunError(err); ok {
				return sr, err
			}
			return singleChunkReader(fmt.Sprintf(m.errorFormat, err)), nil
		}
		return safeWrapReader(sr, m.errorFormat), nil
	}, nil
}

func singleChunkReader(msg string) *schema.StreamReader[string] {
	return schema.StreamReaderFromArray([]string{msg})
}

func safeWrapReader(sr *schema.StreamReader[string], errorFormat string) *schema.StreamReader[string] {
	if sr == nil {
		return schema.StreamReaderFromArray([]string{})
	}
	out, inw := schema.Pipe[string](4)
	go func() {
		defer inw.Close()
		defer sr.Close()
		for {
			chunk, err := sr.Recv()
			if err != nil {
				if errors.Is(err, io.EOF) {
					return
				}
				if _, ok := compose.IsInterruptRerunError(err); ok {
					_ = inw.Send("", err)
					return
				}
				if inw.Send(fmt.Sprintf(errorFormat, err), nil) {
					return
				}
				return
			}
			if inw.Send(chunk, nil) {
				return
			}
		}
	}()
	return out
}
