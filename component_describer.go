package einolib

import (
	"fmt"
	"strings"

	"github.com/cloudwego/eino/components"
)

// 自定义组件类别
const (
	ComponentOfADKAgent      components.Component = "ADKAgent"
	ComponentOfADKOperator   components.Component = "ADKOperator"
	ComponentOfADKMiddleware components.Component = "ADKChatModelAgentMiddleware"
)

// 通用组件名称
const ComponentNameGeneral = "*"

type componentDescriber struct {
	component   components.Component // 组件类别，必需
	Type        string               // 组件类型，必需（不同类别组件下类型可能相同）
	Name        string               // 组件名称，可选（有些组件类型不需要名称，如中间件）
	Description string               // 组件描述信息，可选（不用于注册表中唯一标识）
}

func NewComponentDescriber(component components.Component, typ, name string) componentDescriber {
	return componentDescriber{component: component, Type: typ, Name: name}
}

func (d *componentDescriber) String() string {
	if d == nil {
		return ""
	}
	return fmt.Sprintf("component=%s type=%s name=%s description=%s", d.component, d.Type, d.Name, d.Description)
}

func (d *componentDescriber) Key() string {
	if d == nil {
		return ""
	}
	return strings.ToLower(fmt.Sprintf("%s:%s:%s", d.component, d.Type, d.Name))
}

func (d *componentDescriber) Validate() error {
	if d == nil {
		return fmt.Errorf("describer is nil")
	}
	if d.component == "" {
		return fmt.Errorf("component invalid: empty")
	}
	if d.Type == "" {
		return fmt.Errorf("componentType invalid: empty")
	}
	return nil
}
