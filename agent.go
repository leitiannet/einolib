package einolib

import (
	"context"

	"github.com/cloudwego/eino/adk"
)

// https://www.cloudwego.io/zh/docs/eino/core_modules/eino_adk/
// 智能体类型
type AgentType string

const (
	AgentTypeChatModel          AgentType = "chatmodel"
	AgentTypeDeep               AgentType = "deep"
	AgentTypePlanExecute        AgentType = "planexecute"
	AgentTypeSupervisor         AgentType = "supervisor"
	AgentTypeWorkflowLoop       AgentType = "workflow_loop"
	AgentTypeWorkflowParallel   AgentType = "workflow_parallel"
	AgentTypeWorkflowSequential AgentType = "workflow_sequential"
)

const AgentNameGeneral = ComponentNameGeneral // 通用智能体名称

// 智能体描述
type AgentDescriber struct {
	componentDescriber
}

func NewAgentDescriber(agentType AgentType, agentName string) *AgentDescriber {
	return &AgentDescriber{
		componentDescriber: NewComponentDescriber(ComponentOfADKAgent, string(agentType), agentName),
	}
}

// 智能体配置
type AgentConfig struct {
	componentConfig // 公共元数据
}

// 将名称和描述填充到目标字段
func (ac *AgentConfig) ApplyNameAndDescription(name, description *string) {
	if ac == nil {
		return
	}
	if name != nil && *name == "" {
		*name = ac.Name
	}
	if description != nil && *description == "" {
		*description = ac.Description
	}
}

func NewAgentConfig(agentOptions ...AgentOption) *AgentConfig {
	agentConfig := &AgentConfig{
		componentConfig: NewComponentConfig(ComponentOfADKAgent, "", ""),
	}
	ApplyOptions(agentConfig, agentOptions)
	return agentConfig
}

// 智能体选项
type AgentOption func(agentConfig *AgentConfig)

var (
	WithAgentType            = MakeOption(func(c *AgentConfig, v AgentType) { c.Type = string(v) })
	WithAgentName            = MakeOption(func(c *AgentConfig, v string) { c.Name = v })
	WithAgentDescription     = MakeOption(func(c *AgentConfig, v string) { c.Description = v })
	WithAgentComponentConfig = MakeOption(func(c *AgentConfig, value interface{}) {
		desc, err := LookupAgentDescriber(value)
		if err != nil {
			logger.Warnf("LookupAgentDescriber failed: %v", err)
			return
		}
		if desc == nil {
			logger.Warnf("describer is nil for type %T", value)
			return
		}
		c.SetConfig(desc, value)
	})
)

type AgentConstructor interface {
	Construct(ctx context.Context, agentConfig *AgentConfig) (adk.Agent, error)
}

type AgentConstructFunc func(ctx context.Context, agentConfig *AgentConfig, specificConfig interface{}) (adk.Agent, error)

// 智能体构造器注册中心（类型+名称唯一，大小写无感）
var agentConstructorRegistry = NewComponentRegistry[*AgentDescriber, AgentConstructor]()

func GetAgentConstructor(agentDesc *AgentDescriber) (AgentConstructor, error) {
	return agentConstructorRegistry.GetWithFallback(
		agentDesc,
		AgentNameGeneral,
		func(d *AgentDescriber) string { return d.Name },
		func(d *AgentDescriber) *AgentDescriber {
			return NewAgentDescriber(AgentType(d.Type), AgentNameGeneral)
		},
	)
}

func RegisterAgentConstructor(agentDesc *AgentDescriber, agentConstructor AgentConstructor, bindValues ...interface{}) error {
	return agentConstructorRegistry.Register(agentDesc, agentConstructor, bindValues...)
}

// 注册智能体构造函数；无特定配置时不要传参
func RegisterAgentConstructFunc(agentType AgentType, agentName string, agentConstructFunc AgentConstructFunc, bindValues ...interface{}) error {
	agentDesc := NewAgentDescriber(agentType, agentName)
	agentConstructor := NewComponentConstructor[*AgentDescriber, *AgentConfig, adk.Agent](agentDesc, agentConstructFunc)
	return RegisterAgentConstructor(agentDesc, agentConstructor, bindValues...)
}

func LookupAgentDescriber(value interface{}) (*AgentDescriber, error) {
	return agentConstructorRegistry.LookupDesc(value)
}

func NewAgent(ctx context.Context, agentOptions ...AgentOption) (adk.Agent, error) {
	agentConfig := NewAgentConfig(agentOptions...)
	agentConstructor, err := GetAgentConstructor(NewAgentDescriber(AgentType(agentConfig.Type), agentConfig.Name))
	if err != nil {
		return nil, err
	}
	return agentConstructor.Construct(ctx, agentConfig)
}

func MustNewAgent(ctx context.Context, agentOptions ...AgentOption) adk.Agent {
	agent, err := NewAgent(ctx, agentOptions...)
	if err != nil {
		panic(err)
	}
	if agent == nil {
		panic("MustNewAgent failed: instance is nil")
	}
	return agent
}
