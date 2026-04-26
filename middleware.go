package einolib

import (
	"context"

	"github.com/cloudwego/eino/adk"
)

// 中间件类型
type MiddlewareType string

const (
	MiddlewareTypeAgentsMD       MiddlewareType = "agentsmd"
	MiddlewareTypeFileSystem     MiddlewareType = "filesystem"
	MiddlewareTypePatchToolCalls MiddlewareType = "patchtoolcalls"
	MiddlewareTypePlanTask       MiddlewareType = "plantask"
	MiddlewareTypeReduction      MiddlewareType = "reduction"
	MiddlewareTypeSafeTool       MiddlewareType = "safetool"
	MiddlewareTypeSkill          MiddlewareType = "skill"
	MiddlewareTypeSummarization  MiddlewareType = "summarization"
	MiddlewareTypeToolSearch     MiddlewareType = "toolsearch"
	MiddlewareTypeTrace          MiddlewareType = "trace"
)

const MiddlewareNameGeneral = ComponentNameGeneral // 通用中间件名称

// 中间件描述
type MiddlewareDescriber struct {
	componentDescriber
}

func NewMiddlewareDescriber(middlewareType MiddlewareType, middlewareName string) *MiddlewareDescriber {
	return &MiddlewareDescriber{
		componentDescriber: NewComponentDescriber(ComponentOfADKMiddleware, string(middlewareType), middlewareName),
	}
}

// 中间件配置
type MiddlewareConfig struct {
	componentConfig // 公共元数据
}

func NewMiddlewareConfig(middlewareOptions ...MiddlewareOption) *MiddlewareConfig {
	middlewareConfig := &MiddlewareConfig{
		componentConfig: NewComponentConfig(ComponentOfADKMiddleware, "", ""),
	}
	ApplyOptions(middlewareConfig, middlewareOptions)
	return middlewareConfig
}

// 中间件选项
type MiddlewareOption func(middlewareConfig *MiddlewareConfig)

var (
	WithMiddlewareType            = MakeOption(func(c *MiddlewareConfig, v MiddlewareType) { c.Type = string(v) })
	WithMiddlewareName            = MakeOption(func(c *MiddlewareConfig, v string) { c.Name = v })
	WithMiddlewareDescription     = MakeOption(func(c *MiddlewareConfig, v string) { c.Description = v })
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
	return middlewareConstructorRegistry.GetWithFallback(
		middlewareDesc,
		MiddlewareNameGeneral,
		func(d *MiddlewareDescriber) string { return d.Name },
		func(d *MiddlewareDescriber) *MiddlewareDescriber {
			return NewMiddlewareDescriber(MiddlewareType(d.Type), MiddlewareNameGeneral)
		},
	)
}

func RegisterMiddlewareConstructor(middlewareDesc *MiddlewareDescriber, middlewareConstructor MiddlewareConstructor, bindValues ...interface{}) error {
	return middlewareConstructorRegistry.Register(middlewareDesc, middlewareConstructor, bindValues...)
}

// 注册中间件构造函数；无特定配置时不要传参
func RegisterMiddlewareConstructFunc(middlewareType MiddlewareType, middlewareName string, middlewareConstructFunc MiddlewareConstructFunc, bindValues ...interface{}) error {
	middlewareDesc := NewMiddlewareDescriber(middlewareType, middlewareName)
	middlewareConstructor := NewComponentConstructor[*MiddlewareDescriber, *MiddlewareConfig, adk.ChatModelAgentMiddleware](middlewareDesc, middlewareConstructFunc)
	return RegisterMiddlewareConstructor(middlewareDesc, middlewareConstructor, bindValues...)
}

func LookupMiddlewareDescriber(value interface{}) (*MiddlewareDescriber, error) {
	return middlewareConstructorRegistry.LookupDesc(value)
}

func NewMiddleware(ctx context.Context, middlewareOptions ...MiddlewareOption) (adk.ChatModelAgentMiddleware, error) {
	middlewareConfig := NewMiddlewareConfig(middlewareOptions...)
	middlewareConstructor, err := GetMiddlewareConstructor(NewMiddlewareDescriber(MiddlewareType(middlewareConfig.Type), middlewareConfig.Name))
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
