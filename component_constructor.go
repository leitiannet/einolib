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
	// Go泛型不能直接对约束类型做 == nil 比较，需转为any再判断
	if any(componentConfig) == nil {
		return zero, fmt.Errorf("componentConfig is nil")
	}
	if constructor == nil || constructor.constructFunc == nil {
		return zero, fmt.Errorf("constructor or constructor.constructFunc is nil")
	}
	// 可能为nil
	specificConfig := componentConfig.GetConfig(constructor.desc)
	return constructor.constructFunc(ctx, componentConfig, specificConfig)
}

// 解析特定配置
func ParseSpecificConfig[T any](specificConfig interface{}, defaultConfig func() *T) (*T, error) {
	if specificConfig == nil {
		if defaultConfig == nil {
			return nil, nil
		}
		return defaultConfig(), nil
	}
	v, ok := specificConfig.(*T)
	if !ok || v == nil {
		var want *T
		return nil, fmt.Errorf("specificConfig type mismatch: %T is not %T", specificConfig, want)
	}
	return v, nil
}
