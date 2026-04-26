package einolib

import (
	"fmt"
	"reflect"
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
	bindMap  map[string]D // 绑定类型信息，键为规范化的reflect.Type.String()
	mu       sync.RWMutex // 读写锁
}

func NewComponentRegistry[D ComponentDescriber, V any]() *ComponentRegistry[D, V] {
	return &ComponentRegistry[D, V]{
		valueMap: make(map[string]V),
		bindMap:  make(map[string]D),
	}
}

// 获取值
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

// 获取值(先按具体名、再回退到通用名)
func (r *ComponentRegistry[D, V]) GetWithFallback(
	desc D,
	generalName string,
	getNameFunc func(D) string,
	getGeneralDescFunc func(D) D,
) (V, error) {
	v, err := r.Get(desc)
	if err == nil {
		return v, nil
	}
	if getNameFunc(desc) == generalName {
		return v, err
	}
	fb, err2 := r.Get(getGeneralDescFunc(desc))
	if err2 != nil {
		return v, err
	}
	return fb, nil
}

// 注册组件
func (r *ComponentRegistry[D, V]) Register(desc D, value V, bindTypes ...interface{}) error {
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
	if len(bindTypes) > 1 {
		return fmt.Errorf("at most one bind value, got %d", len(bindTypes))
	}
	key := desc.Key()

	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.valueMap[key]; exists {
		return fmt.Errorf("value already exists: %s", desc)
	}
	r.valueMap[key] = value

	if len(bindTypes) == 0 || bindTypes[0] == nil {
		return nil
	}
	lookupKey, err := bindKeyFrom(bindTypes[0])
	if err != nil {
		delete(r.valueMap, key)
		return err
	}
	descKey := key
	if old, exists := r.bindMap[lookupKey]; exists {
		if any(old) != nil && old.Key() == descKey {
			return nil
		}
		delete(r.valueMap, key)
		return fmt.Errorf("type %q already bound to %s", lookupKey, old)
	}
	r.bindMap[lookupKey] = desc
	return nil
}

// 获取绑定键
func bindKeyFrom(v interface{}) (string, error) {
	if v == nil {
		return "", fmt.Errorf("value is nil")
	}
	var rt reflect.Type
	switch t := v.(type) {
	case reflect.Type:
		if t == nil {
			return "", fmt.Errorf("reflect type is nil")
		}
		rt = t
	default:
		rt = reflect.TypeOf(v)
	}
	if rt == nil {
		return "", fmt.Errorf("invalid reflect type")
	}
	if rt.Kind() != reflect.Ptr {
		rt = reflect.PointerTo(rt)
	}
	return rt.String(), nil
}

// 查找描述
func (r *ComponentRegistry[D, V]) LookupDesc(v interface{}) (D, error) {
	var zero D
	if v == nil {
		return zero, fmt.Errorf("value is nil")
	}
	lookupKey, err := bindKeyFrom(v)
	if err != nil {
		return zero, err
	}
	return r.getBind(lookupKey)
}

// 获取绑定
func (r *ComponentRegistry[D, V]) getBind(lookupKey string) (D, error) {
	var zero D
	r.mu.RLock()
	defer r.mu.RUnlock()
	desc, ok := r.bindMap[lookupKey]
	if !ok {
		return zero, fmt.Errorf("value not found: %s", lookupKey)
	}
	return desc, nil
}
