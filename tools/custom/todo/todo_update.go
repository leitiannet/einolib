package todo

import (
	"context"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"
	"github.com/leitiannet/einolib"
)

const (
	UpdateTodoToolName = "update_todo"
	UpdateTodoToolDesc = "Update a todo item, eg: content,deadline..."
)

type UpdateTodoParams struct {
	ID        string  `json:"id" jsonschema:"description=id of the todo"`
	Content   *string `json:"content,omitempty" jsonschema:"description=content of the todo"`
	StartedAt *int64  `json:"started_at,omitempty" jsonschema:"description=start time in unix timestamp"`
	Deadline  *int64  `json:"deadline,omitempty" jsonschema:"description=deadline of the todo in unix timestamp"`
	Done      *bool   `json:"done,omitempty" jsonschema:"description=done status"`
}

func UpdateTodoFunc(ctx context.Context, params *UpdateTodoParams) (string, error) {
	// Mock处理逻辑
	return `{"msg": "update todo success"}`, nil
}

func NewUpdateTodoTool(ctx context.Context, toolConfig *einolib.ToolConfig, todoToolConfig *TodoToolConfig) ([]tool.BaseTool, error) {
	t, err := utils.InferTool(UpdateTodoToolName, UpdateTodoToolDesc, UpdateTodoFunc)
	if err != nil {
		return nil, err
	}
	return []tool.BaseTool{t}, nil
}

func createUpdateTodoTool(ctx context.Context, toolConfig *einolib.ToolConfig, specificConfig interface{}) ([]tool.BaseTool, error) {
	todoToolConfig, err := einolib.ParseSpecificConfig(specificConfig, func() *TodoToolConfig { return NewTodoToolConfig() })
	if err != nil {
		return nil, err
	}
	return NewUpdateTodoTool(ctx, toolConfig, todoToolConfig)
}
