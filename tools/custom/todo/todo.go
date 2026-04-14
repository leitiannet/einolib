// Todo 列表示例自定义工具（增删改查）
package todo

import (
	"github.com/leitiannet/einolib"
)

type TodoToolConfig struct{}

func NewTodoToolConfig(todoToolOptions ...TodoToolOption) *TodoToolConfig {
	todoToolConfig := &TodoToolConfig{}
	einolib.ApplyOptions(todoToolConfig, todoToolOptions)
	return todoToolConfig
}

type TodoToolOption func(*TodoToolConfig)

func init() {
	for _, reg := range []struct {
		name string
		fn   einolib.ToolConstructFunc
	}{
		{CoreTodoToolName, createCoreTodoTool},
		{AddTodoToolName, createAddTodoTool},
		{UpdateTodoToolName, createUpdateTodoTool},
		{ListTodoToolName, createListTodoTool},
	} {
		if err := einolib.RegisterToolConstructFunc(einolib.ToolTypeCustom, reg.name, reg.fn, (*TodoToolConfig)(nil)); err != nil {
			einolib.GetLogger().Errorf("register tool %s failed: %v", reg.name, err)
		}
	}
}
