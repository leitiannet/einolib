package todo

import (
	"context"

	"github.com/cloudwego/eino/components/tool"
	"github.com/leitiannet/einolib"
)

const CoreTodoToolName = "core_todo"

func NewCoreTodoTool(ctx context.Context, toolConfig *einolib.ToolConfig, todoToolConfig *TodoToolConfig) ([]tool.BaseTool, error) {
	addTools, err := NewAddTodoTool(ctx, toolConfig, todoToolConfig)
	if err != nil {
		return nil, err
	}
	add := addTools[0]
	updateTools, err := NewUpdateTodoTool(ctx, toolConfig, todoToolConfig)
	if err != nil {
		return nil, err
	}
	upd := updateTools[0]
	listTools, err := NewListTodoTool(ctx, toolConfig, todoToolConfig)
	if err != nil {
		return nil, err
	}
	lst := listTools[0]
	return []tool.BaseTool{add, upd, lst}, nil
}

func createCoreTodoTool(ctx context.Context, toolConfig *einolib.ToolConfig, specificConfig interface{}) ([]tool.BaseTool, error) {
	todoToolConfig, err := einolib.ParseSpecificConfig(specificConfig, func() *TodoToolConfig { return NewTodoToolConfig() })
	if err != nil {
		return nil, err
	}
	return NewCoreTodoTool(ctx, toolConfig, todoToolConfig)
}
