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
