package einolib

import "github.com/cloudwego/eino/components"

// 组件配置接口
type ComponentConfiger interface {
	GetConfig(desc ComponentDescriber) interface{}
	SetConfig(desc ComponentDescriber, value interface{})
}

type componentConfig struct {
	componentDescriber          // 组件描述
	configMap          *SyncMap // 存储特定组件的配置信息
}

func NewComponentConfig(component components.Component, typ, name string) componentConfig {
	config := &componentConfig{
		componentDescriber: NewComponentDescriber(component, typ, name),
		configMap:          NewSyncMap(),
	}
	return *config
}

func (cc *componentConfig) GetConfig(desc ComponentDescriber) interface{} {
	if cc == nil || cc.configMap == nil {
		return nil
	}
	value, ok := cc.configMap.Get(desc.Key())
	if !ok {
		return nil
	}
	return value
}

// nil会作为有效值写入，表示存在特定组件的配置
func (cc *componentConfig) SetConfig(desc ComponentDescriber, value interface{}) {
	if cc == nil {
		return
	}
	if cc.configMap == nil {
		cc.configMap = NewSyncMap()
	}
	cc.configMap.Set(desc.Key(), value)
}
