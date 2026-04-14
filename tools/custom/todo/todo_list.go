package todo

import (
	"context"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/schema"
	"github.com/leitiannet/einolib"
)

const (
	ListTodoToolName = "list_todo"
	ListTodoToolDesc = "List all todo items"
)

type ListTodoTool struct{}

func (lt *ListTodoTool) Info(ctx context.Context) (*schema.ToolInfo, error) {
	return &schema.ToolInfo{
		Name: ListTodoToolName,
		Desc: ListTodoToolDesc,
		ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
			"finished": {
				Desc:     "filter todo items if finished",
				Type:     schema.Boolean,
				Required: false,
			},
		}),
	}, nil
}

func (lt *ListTodoTool) InvokableRun(ctx context.Context, argumentsInJSON string, opts ...tool.Option) (string, error) {
	// Mock调用逻辑
	return `{"todos": [{"id": "1", "content": "在2024年12月10日之前完成Eino项目演示文稿的准备工作", "started_at": 1717401600, "deadline": 1717488000, "done": false}]}`, nil
}

func NewListTodoTool(ctx context.Context, toolConfig *einolib.ToolConfig, todoToolConfig *TodoToolConfig) ([]tool.BaseTool, error) {
	return []tool.BaseTool{&ListTodoTool{}}, nil
}

func createListTodoTool(ctx context.Context, toolConfig *einolib.ToolConfig, specificConfig interface{}) ([]tool.BaseTool, error) {
	todoToolConfig, err := einolib.ParseSpecificConfig(specificConfig, func() *TodoToolConfig { return NewTodoToolConfig() })
	if err != nil {
		return nil, err
	}
	return NewListTodoTool(ctx, toolConfig, todoToolConfig)
}
