package todo

import (
	"context"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"
	"github.com/cloudwego/eino/schema"
	"github.com/leitiannet/einolib"
)

const (
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

func getTodoTools(names ...string) ([]tool.BaseTool, error) {
	toolInstances := make([]tool.BaseTool, 0, len(names))
	for _, name := range names {
		var toolInstance tool.BaseTool
		switch name {
		case AddTodoToolName:
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
		case UpdateTodoToolName:
			t, err := utils.InferTool(UpdateTodoToolName, UpdateTodoToolDesc, UpdateTodoFunc)
			if err != nil {
				return nil, err
			}
			toolInstance = t
		case ListTodoToolName:
			toolInstance = &ListTodoTool{}
		}
		if toolInstance != nil {
			toolInstances = append(toolInstances, toolInstance)
		}
	}
	return toolInstances, nil
}

func NewCoreTodoTool(_ context.Context, _ *einolib.ToolConfig, _ interface{}) ([]tool.BaseTool, error) {
	return getTodoTools(AddTodoToolName, UpdateTodoToolName, ListTodoToolName)
}

func NewAddTodoTool(_ context.Context, _ *einolib.ToolConfig, _ interface{}) ([]tool.BaseTool, error) {
	return getTodoTools(AddTodoToolName)
}

func NewUpdateTodoTool(_ context.Context, _ *einolib.ToolConfig, _ interface{}) ([]tool.BaseTool, error) {
	return getTodoTools(UpdateTodoToolName)
}

func NewListTodoTool(_ context.Context, _ *einolib.ToolConfig, _ interface{}) ([]tool.BaseTool, error) {
	return getTodoTools(ListTodoToolName)
}

func init() {
	for _, reg := range []struct {
		name string
		fn   einolib.ToolConstructFunc
	}{
		{CoreTodoToolName, NewCoreTodoTool},
		{AddTodoToolName, NewAddTodoTool},
		{UpdateTodoToolName, NewUpdateTodoTool},
		{ListTodoToolName, NewListTodoTool},
	} {
		if err := einolib.RegisterToolConstructFunc(einolib.ToolTypeCustom, reg.name, reg.fn); err != nil {
			einolib.GetLogger().Errorf("register tool %s failed: %v", reg.name, err)
		}
	}
}
