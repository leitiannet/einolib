package einolib

import (
	"context"
	"fmt"
	"strings"

	"github.com/cloudwego/eino/adk"
)

// 中间件类型
type MiddlewareType string

// 其他中间件类型由middlewares子包定义
const MiddlewareTypeUnknown MiddlewareType = "" // 未知类型中间件

const GeneralMiddlewareName = "*" // 通用中间件名称

// 中间件描述
type MiddlewareDescriber struct {
	MiddlewareType MiddlewareType // 中间件类型
	MiddlewareName string         // 中间件名称
}

func NewMiddlewareDescriber(middlewareType MiddlewareType, middlewareName string) *MiddlewareDescriber {
	return &MiddlewareDescriber{
		MiddlewareType: middlewareType,
		MiddlewareName: middlewareName,
	}
}

func (md *MiddlewareDescriber) String() string {
	return fmt.Sprintf("%s:%s", md.MiddlewareType, md.MiddlewareName)
}

func (md *MiddlewareDescriber) Key() string {
	return strings.ToLower(md.String())
}

func (md *MiddlewareDescriber) Validate() error {
	if md.MiddlewareType == MiddlewareTypeUnknown {
		return fmt.Errorf("middlewareType invalid: %q", md.MiddlewareType)
	}
	if md.MiddlewareName == "" {
		return fmt.Errorf("middlewareName invalid: %q", md.MiddlewareName)
	}
	return nil
}

// 中间件配置
type MiddlewareConfig struct {
	ComponentConfig                      // 特定配置
	MiddlewareType        MiddlewareType // 中间件类型
	MiddlewareName        string         // 中间件名称
	MiddlewareDescription string         // 中间件描述
}

func NewMiddlewareConfig(middlewareOptions ...MiddlewareOption) *MiddlewareConfig {
	middlewareConfig := &MiddlewareConfig{
		ComponentConfig: ComponentConfig{
			ConfigMap: NewSyncMap(),
		},
	}
	ApplyOptions(middlewareConfig, middlewareOptions)
	return middlewareConfig
}

// 中间件选项
type MiddlewareOption func(middlewareConfig *MiddlewareConfig)

var (
	WithMiddlewareType            = MakeOption(func(c *MiddlewareConfig, v MiddlewareType) { c.MiddlewareType = v })
	WithMiddlewareName            = MakeOption(func(c *MiddlewareConfig, v string) { c.MiddlewareName = v })
	WithMiddlewareDescription     = MakeOption(func(c *MiddlewareConfig, v string) { c.MiddlewareDescription = v })
	WithMiddlewareComponentConfig = MakeOption(func(c *MiddlewareConfig, value interface{}) {
		desc, err := LookupMiddlewareDescriber(value)
		if err != nil {
			logger.Warnf("LookupMiddlewareDescriber failed: %v", err)
			return
		}
		if desc == nil {
			logger.Warnf("describer is nil for type %T", value)
			return
		}
		c.SetConfig(desc, value)
	})
)

type MiddlewareConstructor interface {
	Construct(ctx context.Context, middlewareConfig *MiddlewareConfig) (adk.ChatModelAgentMiddleware, error)
}

type MiddlewareConstructFunc func(ctx context.Context, middlewareConfig *MiddlewareConfig, specificConfig interface{}) (adk.ChatModelAgentMiddleware, error)

// 中间件构造器注册中心（类型+名称唯一，大小写无感）
var middlewareConstructorRegistry = NewComponentRegistry[*MiddlewareDescriber, MiddlewareConstructor]()

func GetMiddlewareConstructor(middlewareDesc *MiddlewareDescriber) (MiddlewareConstructor, error) {
	constructor, err := middlewareConstructorRegistry.Get(middlewareDesc)
	if err != nil && middlewareDesc.MiddlewareName != GeneralMiddlewareName {
		generalConstructor, generalErr := middlewareConstructorRegistry.Get(NewMiddlewareDescriber(middlewareDesc.MiddlewareType, GeneralMiddlewareName))
		if generalErr != nil {
			return constructor, err
		}
		return generalConstructor, nil
	}
	return constructor, err
}

func RegisterMiddlewareConstructor(middlewareDesc *MiddlewareDescriber, middlewareConstructor MiddlewareConstructor) error {
	return middlewareConstructorRegistry.Register(middlewareDesc, middlewareConstructor)
}

// 注册中间件构造函数；无特定配置时不要传参
func RegisterMiddlewareConstructFunc(middlewareType MiddlewareType, middlewareName string, middlewareConstructFunc MiddlewareConstructFunc, bindValues ...interface{}) error {
	middlewareDesc := NewMiddlewareDescriber(middlewareType, middlewareName)
	middlewareConstructor := NewComponentConstructor[*MiddlewareDescriber, *MiddlewareConfig, adk.ChatModelAgentMiddleware](middlewareDesc, middlewareConstructFunc)
	return middlewareConstructorRegistry.Register(middlewareDesc, middlewareConstructor, bindValues...)
}

func LookupMiddlewareDescriber(value interface{}) (*MiddlewareDescriber, error) {
	return middlewareConstructorRegistry.LookupDesc(value)
}

func NewMiddleware(ctx context.Context, middlewareOptions ...MiddlewareOption) (adk.ChatModelAgentMiddleware, error) {
	middlewareConfig := NewMiddlewareConfig(middlewareOptions...)
	middlewareConstructor, err := GetMiddlewareConstructor(NewMiddlewareDescriber(middlewareConfig.MiddlewareType, middlewareConfig.MiddlewareName))
	if err != nil {
		return nil, err
	}
	return middlewareConstructor.Construct(ctx, middlewareConfig)
}

func MustNewMiddleware(ctx context.Context, middlewareOptions ...MiddlewareOption) adk.ChatModelAgentMiddleware {
	middleware, err := NewMiddleware(ctx, middlewareOptions...)
	if err != nil {
		panic(err)
	}
	if middleware == nil {
		panic("MustNewMiddleware failed: instance is nil")
	}
	return middleware
}
