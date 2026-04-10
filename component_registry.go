package einolib

import (
	"fmt"
	"sync"
)

// 组件描述接口
type ComponentDescriber interface {
	String() string // 字符串描述（用于日志打印）
	Key() string    // 唯一标识符（用于注册查找）
}

// 组件验证接口
type ComponentValidator interface {
	Validate() error // 合法性验证
}

// 通用组件注册器
type ComponentRegistry[D ComponentDescriber, V any] struct {
	valueMap map[string]V // 存储值信息，键为ComponentDescriber.Key()
	mu       sync.RWMutex // 读写锁
}

func NewComponentRegistry[D ComponentDescriber, V any]() *ComponentRegistry[D, V] {
	return &ComponentRegistry[D, V]{
		valueMap: make(map[string]V),
	}
}

// 根据描述获取值
func (r *ComponentRegistry[D, V]) Get(desc D) (V, error) {
	var zero V
	if any(desc) == nil {
		return zero, fmt.Errorf("desc is nil")
	}
	if validator, ok := any(desc).(ComponentValidator); ok {
		if err := validator.Validate(); err != nil {
			return zero, fmt.Errorf("validate failed: %v", err)
		}
	}
	key := desc.Key()

	r.mu.RLock()
	defer r.mu.RUnlock()
	value, exists := r.valueMap[key]
	if !exists {
		return zero, fmt.Errorf("value not found: %s", desc)
	}
	return value, nil
}

// 注册组件
func (r *ComponentRegistry[D, V]) Register(desc D, value V) error {
	if any(desc) == nil {
		return fmt.Errorf("desc is nil")
	}
	if validator, ok := any(desc).(ComponentValidator); ok {
		if err := validator.Validate(); err != nil {
			return fmt.Errorf("validate failed: %v", err)
		}
	}
	if any(value) == nil {
		return fmt.Errorf("value is nil: %s", desc)
	}
	key := desc.Key()

	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.valueMap[key]; exists {
		return fmt.Errorf("value already exists: %s", desc)
	}
	r.valueMap[key] = value
	return nil
}
