package einolib

import (
	"context"

	"github.com/cloudwego/eino/components"
	"github.com/cloudwego/eino/components/model"
	"github.com/kelseyhightower/envconfig"
)

// https://www.cloudwego.io/zh/docs/eino/ecosystem_integration/chat_model/
// 模型类型
type ModelType string

const (
	ModelTypeARK    ModelType = "ark"
	ModelTypeOllama ModelType = "ollama"
	ModelTypeOpenAI ModelType = "openai"
)

const ModelNameGeneral = ComponentNameGeneral // 通用模型名称

// 模型描述
type ModelDescriber struct {
	componentDescriber
}

func NewModelDescriber(modelType ModelType) *ModelDescriber {
	return &ModelDescriber{
		componentDescriber: NewComponentDescriber(components.ComponentOfChatModel, string(modelType), ""),
	}
}

// 模型配置
type ModelConfig struct {
	componentConfig `json:"-" envconfig:"-"` // 公共元数据
	ID              string                   `json:"id" envconfig:"MODEL"`          // 模型ID
	BaseURL         string                   `json:"base_url" envconfig:"BASE_URL"` // 服务地址
	APIKey          string                   `json:"api_key" envconfig:"API_KEY"`   // API密钥
	ByAzure         string                   `json:"by_azure" envconfig:"BY_AZURE"` // 是否使用Azure OpenAI服务
}

func NewModelConfig(modelOptions ...ModelOption) *ModelConfig {
	modelConfig := &ModelConfig{
		componentConfig: NewComponentConfig(components.ComponentOfChatModel, "", ""),
	}
	modelConfig.initFromEnvironment()
	ApplyOptions(modelConfig, modelOptions)
	return modelConfig
}

// 从环境变量初始化配置，支持通用前缀和模型前缀，优先级为通用前缀覆盖
// 通用前缀 EINO_: EINO_MODEL_TYPE, EINO_MODEL, EINO_BASE_URL, EINO_API_KEY, EINO_BY_AZURE
// 模型前缀 {MODEL_TYPE}_: MODEL_TYPE、{MODEL_TYPE}_MODEL, {MODEL_TYPE}_BASE_URL, {MODEL_TYPE}_API_KEY, {MODEL_TYPE}_BY_AZURE
func (modelConfig *ModelConfig) initFromEnvironment() {
	if modelConfig == nil {
		return
	}
	if err := envconfig.Process("EINO", modelConfig); err != nil {
		logger.Warnf("load environment variables failed: %v", err)
	}
	// 不同模型使用不同环境变量前缀
	BindVarFromEnv(&modelConfig.Type, "MODEL_TYPE", "EINO")
	BindVarFromEnv(&modelConfig.Type, "MODEL_TYPE")
	if modelConfig.Type != "" {
		prefix := modelConfig.Type
		BindVarFromEnv(&modelConfig.ID, "MODEL", prefix)
		BindVarFromEnv(&modelConfig.BaseURL, "BASE_URL", prefix)
		BindVarFromEnv(&modelConfig.APIKey, "API_KEY", prefix)
		BindVarFromEnv(&modelConfig.ByAzure, "BY_AZURE", prefix)
	}
}

// 模型选项
type ModelOption func(modelConfig *ModelConfig)

var (
	WithModelType        = MakeOption(func(c *ModelConfig, v ModelType) { c.Type = string(v) })
	WithModelName        = MakeOption(func(c *ModelConfig, v string) { c.Name = v })
	WithModelDescription = MakeOption(func(c *ModelConfig, v string) { c.Description = v })
	WithModelID          = MakeOption(func(c *ModelConfig, v string) { c.ID = v })
	WithModelBaseURL     = MakeOption(func(c *ModelConfig, v string) { c.BaseURL = v })
	WithModelAPIKey      = MakeOption(func(c *ModelConfig, v string) { c.APIKey = v })
	WithModelByAzure     = MakeOption(func(c *ModelConfig, v string) { c.ByAzure = v })
	WithModelByAzureBool = MakeOption(func(c *ModelConfig, v bool) {
		if v {
			c.ByAzure = "true"
		} else {
			c.ByAzure = "false"
		}
	})
	WithModelComponentConfig = MakeOption(func(c *ModelConfig, value interface{}) {
		desc, err := LookupModelDescriber(value)
		if err != nil {
			logger.Warnf("LookupModelDescriber failed: %v", err)
			return
		}
		if desc == nil {
			logger.Warnf("describer is nil for type %T", value)
			return
		}
		c.SetConfig(desc, value)
	})
)

type ModelConstructor interface {
	Construct(ctx context.Context, modelConfig *ModelConfig) (model.ToolCallingChatModel, error)
}

type ModelConstructFunc func(ctx context.Context, modelConfig *ModelConfig, specificConfig interface{}) (model.ToolCallingChatModel, error)

// 模型构造器注册中心（类型唯一，大小写无感）
var modelConstructorRegistry = NewComponentRegistry[*ModelDescriber, ModelConstructor]()

func GetModelConstructor(modelDesc *ModelDescriber) (ModelConstructor, error) {
	return modelConstructorRegistry.Get(modelDesc)
}

func RegisterModelConstructor(modelDesc *ModelDescriber, modelConstructor ModelConstructor, bindValues ...interface{}) error {
	return modelConstructorRegistry.Register(modelDesc, modelConstructor, bindValues...)
}

// 注册模型构造函数；无特定配置时不要传参
func RegisterModelConstructFunc(modelType ModelType, modelConstructFunc ModelConstructFunc, bindValues ...interface{}) error {
	modelDesc := NewModelDescriber(modelType)
	modelConstructor := NewComponentConstructor[*ModelDescriber, *ModelConfig, model.ToolCallingChatModel](modelDesc, modelConstructFunc)
	return RegisterModelConstructor(modelDesc, modelConstructor, bindValues...)
}

func LookupModelDescriber(value interface{}) (*ModelDescriber, error) {
	return modelConstructorRegistry.LookupDesc(value)
}

func NewChatModel(ctx context.Context, modelOptions ...ModelOption) (model.ToolCallingChatModel, error) {
	modelConfig := NewModelConfig(modelOptions...)
	modelConstructor, err := GetModelConstructor(NewModelDescriber(ModelType(modelConfig.Type)))
	if err != nil {
		return nil, err
	}
	return modelConstructor.Construct(ctx, modelConfig)
}

func MustNewChatModel(ctx context.Context, modelOptions ...ModelOption) model.ToolCallingChatModel {
	chatModel, err := NewChatModel(ctx, modelOptions...)
	if err != nil {
		panic(err)
	}
	if chatModel == nil {
		panic("MustNewChatModel failed: instance is nil")
	}
	return chatModel
}
