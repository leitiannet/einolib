package einolib

import (
	"context"

	"github.com/cloudwego/eino/adk/filesystem"
)

// 执行器接口
type Operator interface {
	filesystem.Backend
	filesystem.Shell
	filesystem.StreamingShell
}

// 执行器类型
type OperatorType string

const (
	OperatorTypeAgentKit OperatorType = "agentkit"
	OperatorTypeLocal    OperatorType = "local"
	OperatorTypeMemory   OperatorType = "memory"
)

const OperatorNameGeneral = ComponentNameGeneral // 通用执行器名称

// 执行器描述
type OperatorDescriber struct {
	componentDescriber
}

func NewOperatorDescriber(operatorType OperatorType, operatorName string) *OperatorDescriber {
	return &OperatorDescriber{
		componentDescriber: NewComponentDescriber(ComponentOfADKOperator, string(operatorType), operatorName),
	}
}

// 执行器配置
type OperatorConfig struct {
	componentConfig // 公共元数据
}

func NewOperatorConfig(operatorOptions ...OperatorOption) *OperatorConfig {
	operatorConfig := &OperatorConfig{
		componentConfig: NewComponentConfig(ComponentOfADKOperator, "", ""),
	}
	ApplyOptions(operatorConfig, operatorOptions)
	return operatorConfig
}

// 执行器选项
type OperatorOption func(operatorConfig *OperatorConfig)

var (
	WithOperatorType            = MakeOption(func(c *OperatorConfig, v OperatorType) { c.Type = string(v) })
	WithOperatorName            = MakeOption(func(c *OperatorConfig, v string) { c.Name = v })
	WithOperatorDescription     = MakeOption(func(c *OperatorConfig, v string) { c.Description = v })
	WithOperatorComponentConfig = MakeOption(func(c *OperatorConfig, value interface{}) {
		desc, err := LookupOperatorDescriber(value)
		if err != nil {
			logger.Warnf("LookupOperatorDescriber failed: %v", err)
			return
		}
		if desc == nil {
			logger.Warnf("describer is nil for type %T", value)
			return
		}
		c.SetConfig(desc, value)
	})
)

type OperatorConstructor interface {
	Construct(ctx context.Context, operatorConfig *OperatorConfig) (Operator, error)
}

type OperatorConstructFunc func(ctx context.Context, operatorConfig *OperatorConfig, specificConfig interface{}) (Operator, error)

// 执行器构造器注册中心（类型+名称唯一，大小写无感）
var operatorConstructorRegistry = NewComponentRegistry[*OperatorDescriber, OperatorConstructor]()

func GetOperatorConstructor(operatorDesc *OperatorDescriber) (OperatorConstructor, error) {
	return operatorConstructorRegistry.GetWithFallback(
		operatorDesc,
		OperatorNameGeneral,
		func(d *OperatorDescriber) string { return d.Name },
		func(d *OperatorDescriber) *OperatorDescriber {
			return NewOperatorDescriber(OperatorType(d.Type), OperatorNameGeneral)
		},
	)
}

func RegisterOperatorConstructor(operatorDesc *OperatorDescriber, operatorConstructor OperatorConstructor, bindValues ...interface{}) error {
	return operatorConstructorRegistry.Register(operatorDesc, operatorConstructor, bindValues...)
}

// 注册执行器构造函数；无特定配置时不要传参
func RegisterOperatorConstructFunc(operatorType OperatorType, operatorName string, operatorConstructFunc OperatorConstructFunc, bindValues ...interface{}) error {
	operatorDesc := NewOperatorDescriber(operatorType, operatorName)
	operatorConstructor := NewComponentConstructor[*OperatorDescriber, *OperatorConfig, Operator](operatorDesc, operatorConstructFunc)
	return RegisterOperatorConstructor(operatorDesc, operatorConstructor, bindValues...)
}

func LookupOperatorDescriber(value interface{}) (*OperatorDescriber, error) {
	return operatorConstructorRegistry.LookupDesc(value)
}

func NewOperator(ctx context.Context, operatorOptions ...OperatorOption) (Operator, error) {
	operatorConfig := NewOperatorConfig(operatorOptions...)
	operatorConstructor, err := GetOperatorConstructor(NewOperatorDescriber(OperatorType(operatorConfig.Type), operatorConfig.Name))
	if err != nil {
		return nil, err
	}
	return operatorConstructor.Construct(ctx, operatorConfig)
}

func MustNewOperator(ctx context.Context, operatorOptions ...OperatorOption) Operator {
	operator, err := NewOperator(ctx, operatorOptions...)
	if err != nil {
		panic(err)
	}
	if operator == nil {
		panic("MustNewOperator failed: instance is nil")
	}
	return operator
}
