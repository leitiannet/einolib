package todo

import (
	"context"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"
	"github.com/cloudwego/eino/schema"
	"github.com/leitiannet/einolib"
)

const (
	// 核心工具集标识，聚合add_todo、update_todo、list_todo
	CoreTodoToolName   = "core_todo"
	AddTodoToolName    = "add_todo"
	AddTodoToolDesc    = "Add a todo item"
	UpdateTodoToolName = "update_todo"
	UpdateTodoToolDesc = "Update a todo item, eg: content,deadline..."
	ListTodoToolName   = "list_todo"
	ListTodoToolDesc   = "List all todo items"
)

type AddTodoParams struct {
	Content  string `json:"content"`              // 具体内容
	StartAt  *int64 `json:"started_at,omitempty"` // 开始时间
	Deadline *int64 `json:"deadline,omitempty"`   // 截止时间
}

func AddTodoFunc(_ context.Context, params *AddTodoParams) (string, error) {
	// Mock处理逻辑
	return `{"msg": "add todo success"}`, nil
}

type UpdateTodoParams struct {
	ID        string  `json:"id" jsonschema:"description=id of the todo"`
	Content   *string `json:"content,omitempty" jsonschema:"description=content of the todo"`
	StartedAt *int64  `json:"started_at,omitempty" jsonschema:"description=start time in unix timestamp"`
	Deadline  *int64  `json:"deadline,omitempty" jsonschema:"description=deadline of the todo in unix timestamp"`
	Done      *bool   `json:"done,omitempty" jsonschema:"description=done status"`
}

func UpdateTodoFunc(_ context.Context, params *UpdateTodoParams) (string, error) {
	// Mock处理逻辑
	return `{"msg": "update todo success"}`, nil
}

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

func getTodoTool(names ...string) ([]tool.BaseTool, error) {
	toolInstances := make([]tool.BaseTool, 0)
	for _, name := range names {
		var toolInstance tool.BaseTool
		var err error
		switch name {
		case AddTodoToolName:
			// 使用 NewTool 创建工具
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
			toolInstance = utils.NewTool(info, AddTodoFunc)
			err = nil
		case UpdateTodoToolName:
			// 使用 InferTool 创建工具
			toolInstance, err = utils.InferTool(
				UpdateTodoToolName,
				UpdateTodoToolDesc,
				UpdateTodoFunc)
		case ListTodoToolName:
			toolInstance = &ListTodoTool{}
			err = nil
		}
		if err != nil {
			return nil, err
		}
		if toolInstance != nil {
			toolInstances = append(toolInstances, toolInstance)
		}
	}
	return toolInstances, nil
}

func init() {
	_ = einolib.RegisterToolConstructFunc(einolib.ToolTypeCustom, CoreTodoToolName, func(ctx context.Context, toolConfig *einolib.ToolConfig, specificConfig interface{}) ([]tool.BaseTool, error) {
		return getTodoTool(AddTodoToolName, UpdateTodoToolName, ListTodoToolName)
	})
	_ = einolib.RegisterToolConstructFunc(einolib.ToolTypeCustom, AddTodoToolName, func(ctx context.Context, toolConfig *einolib.ToolConfig, specificConfig interface{}) ([]tool.BaseTool, error) {
		return getTodoTool(AddTodoToolName)
	})
	_ = einolib.RegisterToolConstructFunc(einolib.ToolTypeCustom, UpdateTodoToolName, func(ctx context.Context, toolConfig *einolib.ToolConfig, specificConfig interface{}) ([]tool.BaseTool, error) {
		return getTodoTool(UpdateTodoToolName)
	})
	_ = einolib.RegisterToolConstructFunc(einolib.ToolTypeCustom, ListTodoToolName, func(ctx context.Context, toolConfig *einolib.ToolConfig, specificConfig interface{}) ([]tool.BaseTool, error) {
		return getTodoTool(ListTodoToolName)
	})
}
