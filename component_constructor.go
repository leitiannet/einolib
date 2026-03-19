package einolib

import (
	"context"
	"fmt"
)

// 通用组件构造器
type ComponentConstructor[D ComponentDescriber, C ComponentConfiger, R any] struct {
	desc          D
	constructFunc func(ctx context.Context, componentConfig C, specificConfig interface{}) (R, error)
}

func NewComponentConstructor[D ComponentDescriber, C ComponentConfiger, R any](
	desc D,
	constructFunc func(ctx context.Context, componentConfig C, specificConfig interface{}) (R, error),
) *ComponentConstructor[D, C, R] {
	return &ComponentConstructor[D, C, R]{
		desc:          desc,
		constructFunc: constructFunc,
	}
}

func (constructor *ComponentConstructor[D, C, R]) Construct(ctx context.Context, componentConfig C) (R, error) {
	var zero R
	// 转为interface{}再比较
	var emptyInterface interface{} = componentConfig
	if emptyInterface == nil {
		return zero, fmt.Errorf("componentConfig is nil")
	}
	if constructor == nil || constructor.constructFunc == nil {
		return zero, fmt.Errorf("constructor or constructor.constructFunc is nil")
	}
	// 可能为nil
	specificConfig := componentConfig.GetConfig(constructor.desc)
	return constructor.constructFunc(ctx, componentConfig, specificConfig)
}
