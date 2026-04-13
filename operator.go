package einolib

import (
	"context"
	"fmt"
	"strings"

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

// 其他执行器类型由operators子包定义
const OperatorTypeUnknown OperatorType = "" // 未知类型执行器

const GeneralOperatorName = "*" // 通用执行器名称

// 执行器描述
type OperatorDescriber struct {
	OperatorType OperatorType // 执行器类型
	OperatorName string       // 执行器名称
}

func NewOperatorDescriber(operatorType OperatorType, operatorName string) *OperatorDescriber {
	return &OperatorDescriber{
		OperatorType: operatorType,
		OperatorName: operatorName,
	}
}

func (od *OperatorDescriber) String() string {
	return fmt.Sprintf("%s:%s", od.OperatorType, od.OperatorName)
}

func (od *OperatorDescriber) Key() string {
	return strings.ToLower(od.String())
}

func (od *OperatorDescriber) Validate() error {
	if od.OperatorType == OperatorTypeUnknown {
		return fmt.Errorf("operatorType invalid: %q", od.OperatorType)
	}
	if od.OperatorName == "" {
		return fmt.Errorf("operatorName invalid: %q", od.OperatorName)
	}
	return nil
}

// 执行器配置
type OperatorConfig struct {
	ComponentConfig                  // 特定配置
	OperatorType        OperatorType // 执行器类型
	OperatorName        string       // 执行器名称
	OperatorDescription string       // 执行器描述
}

func NewOperatorConfig(operatorOptions ...OperatorOption) *OperatorConfig {
	operatorConfig := &OperatorConfig{
		ComponentConfig: ComponentConfig{
			ConfigMap: NewSyncMap(),
		},
	}
	ApplyOptions(operatorConfig, operatorOptions)
	return operatorConfig
}

// 执行器选项
type OperatorOption func(operatorConfig *OperatorConfig)

var (
	WithOperatorType            = MakeOption(func(c *OperatorConfig, v OperatorType) { c.OperatorType = v })
	WithOperatorName            = MakeOption(func(c *OperatorConfig, v string) { c.OperatorName = v })
	WithOperatorDescription     = MakeOption(func(c *OperatorConfig, v string) { c.OperatorDescription = v })
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
	constructor, err := operatorConstructorRegistry.Get(operatorDesc)
	if err != nil && operatorDesc.OperatorName != GeneralOperatorName {
		generalConstructor, generalErr := operatorConstructorRegistry.Get(NewOperatorDescriber(operatorDesc.OperatorType, GeneralOperatorName))
		if generalErr != nil {
			return constructor, err
		}
		return generalConstructor, nil
	}
	return constructor, err
}

func RegisterOperatorConstructor(operatorDesc *OperatorDescriber, operatorConstructor OperatorConstructor) error {
	return operatorConstructorRegistry.Register(operatorDesc, operatorConstructor)
}

// 注册执行器构造函数；无特定配置时不要传参
func RegisterOperatorConstructFunc(operatorType OperatorType, operatorName string, operatorConstructFunc OperatorConstructFunc, bindValues ...interface{}) error {
	operatorDesc := NewOperatorDescriber(operatorType, operatorName)
	operatorConstructor := NewComponentConstructor[*OperatorDescriber, *OperatorConfig, Operator](operatorDesc, operatorConstructFunc)
	return operatorConstructorRegistry.Register(operatorDesc, operatorConstructor, bindValues...)
}

func LookupOperatorDescriber(value interface{}) (*OperatorDescriber, error) {
	return operatorConstructorRegistry.LookupDesc(value)
}

func NewOperator(ctx context.Context, operatorOptions ...OperatorOption) (Operator, error) {
	operatorConfig := NewOperatorConfig(operatorOptions...)
	operatorConstructor, err := GetOperatorConstructor(NewOperatorDescriber(operatorConfig.OperatorType, operatorConfig.OperatorName))
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
