package einolib

// 组件配置接口
type ComponentConfiger interface {
	GetConfig(desc ComponentDescriber) interface{}
	SetConfig(desc ComponentDescriber, value interface{})
}

type ComponentConfig struct {
	ConfigMap *SyncMap // 存储特定组件的配置信息
}

func (cc *ComponentConfig) GetConfig(desc ComponentDescriber) interface{} {
	if cc == nil || cc.ConfigMap == nil {
		return nil
	}
	value, ok := cc.ConfigMap.Get(desc.Key())
	if !ok {
		return nil
	}
	return value
}

// nil会作为有效值写入，表示存在特定组件的配置
func (cc *ComponentConfig) SetConfig(desc ComponentDescriber, value interface{}) {
	if cc == nil {
		return
	}
	if cc.ConfigMap == nil {
		cc.ConfigMap = NewSyncMap()
	}
	cc.ConfigMap.Set(desc.Key(), value)
}
