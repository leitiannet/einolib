package einolib

// 对target依次应用opts中非nil的选项函数，如果opt为nil，则不应用
// 无~表示类型本身，有~表示类型本身和以类型为底层类型的所有类型
func ApplyOptions[T any, Opt ~func(*T)](target *T, opts []Opt) {
	if target == nil {
		return
	}
	for _, opt := range opts {
		if opt != nil {
			opt(target)
		}
	}
}

func ApplyOptionsVariadic[T any, Opt ~func(*T)](target *T, opts ...Opt) {
	ApplyOptions(target, opts)
}

// 创建接收单个参数的选项函数
func MakeOption[Config any, Value any](setter func(*Config, Value)) func(Value) func(*Config) {
	return func(val Value) func(*Config) {
		return func(cfg *Config) {
			setter(cfg, val)
		}
	}
}

// 创建追加元素到切片字段的选项函数
func MakeAppendOption[Config any, Elem any](getter func(*Config) *[]Elem) func(...Elem) func(*Config) {
	return func(vals ...Elem) func(*Config) {
		return func(cfg *Config) {
			sl := getter(cfg)
			if sl == nil {
				logger.Debugf("target slice pointer is nil")
				return
			}
			if len(vals) == 0 {
				logger.Warnf("no values to append")
				return
			}
			*sl = append(*sl, vals...)
		}
	}
}
