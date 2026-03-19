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
	descMap  map[string]D // 存储描述信息，键为ComponentDescriber.Key()
	valueMap map[string]V // 存储值信息，与descMap一一对应
	mu       sync.RWMutex // 读写锁
}

func NewComponentRegistry[D ComponentDescriber, V any]() *ComponentRegistry[D, V] {
	return &ComponentRegistry[D, V]{
		descMap:  make(map[string]D),
		valueMap: make(map[string]V),
	}
}

// 获取
func (r *ComponentRegistry[D, V]) Get(desc D) (V, error) {
	var zero V
	if validator, ok := any(desc).(ComponentValidator); ok {
		if err := validator.Validate(); err != nil {
			return zero, fmt.Errorf("validate failed: %w", err)
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

// 注册
func (r *ComponentRegistry[D, V]) Register(desc D, value V) error {
	if validator, ok := any(desc).(ComponentValidator); ok {
		if err := validator.Validate(); err != nil {
			return fmt.Errorf("validate failed: %w", err)
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
	r.descMap[key] = desc
	r.valueMap[key] = value
	return nil
}
