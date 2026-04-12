package einolib

import (
	"context"
	"fmt"
	"strings"

	"github.com/cloudwego/eino/adk"
)

// https://www.cloudwego.io/zh/docs/eino/core_modules/eino_adk/
// 智能体类型
type AgentType string

// 其他智能体类型由agents子包定义
const AgentTypeUnknown AgentType = "" // 未知类型智能体

const GeneralAgentName = "*" // 通用智能体名称

// 智能体描述
type AgentDescriber struct {
	AgentType AgentType // 智能体类型
	AgentName string    // 智能体名称
}

func NewAgentDescriber(agentType AgentType, agentName string) *AgentDescriber {
	return &AgentDescriber{
		AgentType: agentType,
		AgentName: agentName,
	}
}

func (ad *AgentDescriber) String() string {
	return fmt.Sprintf("%s:%s", ad.AgentType, ad.AgentName)
}

func (ad *AgentDescriber) Key() string {
	return strings.ToLower(ad.String())
}

func (ad *AgentDescriber) Validate() error {
	if ad.AgentType == AgentTypeUnknown {
		return fmt.Errorf("agentType invalid: %q", ad.AgentType)
	}
	if ad.AgentName == "" {
		return fmt.Errorf("agentName invalid: %q", ad.AgentName)
	}
	return nil
}

// 智能体配置
type AgentConfig struct {
	ComponentConfig            // 特定配置
	AgentType        AgentType // 智能体类型
	AgentName        string    // 智能体名称
	AgentDescription string    // 智能体描述
}

// 将名称和描述填充到目标字段
func (ac *AgentConfig) ApplyNameAndDescription(name, description *string) {
	if ac == nil {
		return
	}
	if name != nil && *name == "" {
		*name = ac.AgentName
	}
	if description != nil && *description == "" {
		*description = ac.AgentDescription
	}
}

func NewAgentConfig(agentOptions ...AgentOption) *AgentConfig {
	agentConfig := &AgentConfig{
		ComponentConfig: ComponentConfig{
			ConfigMap: NewSyncMap(),
		},
	}
	ApplyOptions(agentConfig, agentOptions)
	return agentConfig
}

// 智能体选项
type AgentOption func(agentConfig *AgentConfig)

var (
	WithAgentType            = MakeOption(func(c *AgentConfig, v AgentType) { c.AgentType = v })
	WithAgentName            = MakeOption(func(c *AgentConfig, v string) { c.AgentName = v })
	WithAgentDescription     = MakeOption(func(c *AgentConfig, v string) { c.AgentDescription = v })
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
	constructor, err := agentConstructorRegistry.Get(agentDesc)
	if err != nil && agentDesc.AgentName != GeneralAgentName {
		generalConstructor, generalErr := agentConstructorRegistry.Get(NewAgentDescriber(agentDesc.AgentType, GeneralAgentName))
		if generalErr != nil {
			return constructor, err
		}
		return generalConstructor, nil
	}
	return constructor, err
}

func RegisterAgentConstructor(agentDesc *AgentDescriber, agentConstructor AgentConstructor) error {
	return agentConstructorRegistry.Register(agentDesc, agentConstructor)
}

// 注册智能体构造函数；无特定配置时不要传参
func RegisterAgentConstructFunc(agentType AgentType, agentName string, agentConstructFunc AgentConstructFunc, bindValues ...interface{}) error {
	agentDesc := NewAgentDescriber(agentType, agentName)
	agentConstructor := NewComponentConstructor[*AgentDescriber, *AgentConfig, adk.Agent](agentDesc, agentConstructFunc)
	return agentConstructorRegistry.Register(agentDesc, agentConstructor, bindValues...)
}

func LookupAgentDescriber(value interface{}) (*AgentDescriber, error) {
	return agentConstructorRegistry.LookupDesc(value)
}

func NewAgent(ctx context.Context, agentOptions ...AgentOption) (adk.Agent, error) {
	agentConfig := NewAgentConfig(agentOptions...)
	agentConstructor, err := GetAgentConstructor(NewAgentDescriber(agentConfig.AgentType, agentConfig.AgentName))
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
