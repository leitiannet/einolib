package todo

import (
	"context"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"
	"github.com/cloudwego/eino/schema"
	"github.com/leitiannet/einolib"
)

const (
	AddTodoToolName = "add_todo"
	AddTodoToolDesc = "Add a todo item"
)

type AddTodoParams struct {
	Content  string `json:"content"`              // 具体内容
	StartAt  *int64 `json:"started_at,omitempty"` // 开始时间
	Deadline *int64 `json:"deadline,omitempty"`   // 截止时间
}

func AddTodoFunc(ctx context.Context, params *AddTodoParams) (string, error) {
	// Mock处理逻辑
	return `{"msg": "add todo success"}`, nil
}

func NewAddTodoTool(ctx context.Context, toolConfig *einolib.ToolConfig, todoToolConfig *TodoToolConfig) ([]tool.BaseTool, error) {
	info := &schema.ToolInfo{
		Name: AddTodoToolName,
		Desc: AddTodoToolDesc,
		ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
			"content": {
				Desc:     "The content of the todo item",
				Type:     schema.String,
				Required: true,
			},
			"started_at": {
				Desc: "The started time of the todo item, in unix timestamp",
				Type: schema.Integer,
			},
			"deadline": {
				Desc: "The deadline of the todo item, in unix timestamp",
				Type: schema.Integer,
			},
		}),
	}
	return []tool.BaseTool{utils.NewTool(info, AddTodoFunc)}, nil
}

func createAddTodoTool(ctx context.Context, toolConfig *einolib.ToolConfig, specificConfig interface{}) ([]tool.BaseTool, error) {
	todoToolConfig, err := einolib.ParseSpecificConfig(specificConfig, func() *TodoToolConfig { return NewTodoToolConfig() })
	if err != nil {
		return nil, err
	}
	return NewAddTodoTool(ctx, toolConfig, todoToolConfig)
}
