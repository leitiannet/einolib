package einolib

import (
	"context"
	"fmt"
	"strings"

	"github.com/cloudwego/eino/components/model"
	"github.com/kelseyhightower/envconfig"
)

// https://www.cloudwego.io/zh/docs/eino/ecosystem_integration/chat_model/
// 模型类型
type ModelType string

// 其他模型类型由models子包定义
const ModelTypeUnknown ModelType = "" // 未知类型模型

// 模型描述
type ModelDescriber struct {
	ModelType ModelType // 模型类型
}

func NewModelDescriber(modelType ModelType) *ModelDescriber {
	return &ModelDescriber{
		ModelType: modelType,
	}
}

func (md *ModelDescriber) String() string {
	return string(md.ModelType)
}

func (md *ModelDescriber) Key() string {
	return strings.ToLower(md.String())
}

func (md *ModelDescriber) Validate() error {
	if md.ModelType == ModelTypeUnknown {
		return fmt.Errorf("modelType invalid: %q", md.ModelType)
	}
	return nil
}

// 模型配置
type ModelConfig struct {
	ComponentConfig `json:"-" envconfig:"-"` // 特定配置
	ModelType       ModelType                `json:"model_type" envconfig:"MODEL_TYPE"` // 模型类型
	ModelName       string                   `json:"model_name" envconfig:"MODEL_NAME"` // 模型名称
	BaseURL         string                   `json:"base_url" envconfig:"BASE_URL"`     // 服务地址
	APIKey          string                   `json:"api_key" envconfig:"API_KEY"`       // API密钥
	ByAzure         string                   `json:"by_azure" envconfig:"BY_AZURE"`     // 是否使用Azure OpenAI服务
}

func NewModelConfig(modelOptions ...ModelOption) *ModelConfig {
	modelConfig := &ModelConfig{
		ComponentConfig: ComponentConfig{
			ConfigMap: NewSyncMap(),
		},
	}
	modelConfig.initFromEnvironment()
	ApplyOptions(modelConfig, modelOptions)
	return modelConfig
}

// 从环境变量加载配置，支持两级前缀（模型前缀优先级高于通用前缀）：
// 通用前缀 EINO_: EINO_MODEL_TYPE, EINO_MODEL_NAME, EINO_BASE_URL, EINO_API_KEY, EINO_BY_AZURE
// 模型前缀 {MODEL_TYPE}_: {MODEL_TYPE}_MODEL, {MODEL_TYPE}_BASE_URL, {MODEL_TYPE}_API_KEY, {MODEL_TYPE}_BY_AZURE
func (modelConfig *ModelConfig) initFromEnvironment() {
	if modelConfig == nil {
		return
	}
	if err := envconfig.Process("EINO", modelConfig); err != nil {
		logger.Warnf("load environment variables failed: %v", err)
	}
	// 不同模型使用不同环境变量
	BindVarFromEnv((*string)(&modelConfig.ModelType), "MODEL_TYPE")
	if modelConfig.ModelType != ModelTypeUnknown {
		prefix := string(modelConfig.ModelType)
		BindVarFromEnv(&modelConfig.ModelName, "MODEL", prefix)
		BindVarFromEnv(&modelConfig.BaseURL, "BASE_URL", prefix)
		BindVarFromEnv(&modelConfig.APIKey, "API_KEY", prefix)
		BindVarFromEnv(&modelConfig.ByAzure, "BY_AZURE", prefix)
	}
}

// 模型选项
type ModelOption func(modelConfig *ModelConfig)

var (
	WithModelType   = MakeOption(func(c *ModelConfig, v ModelType) { c.ModelType = v })
	WithModelName   = MakeOption(func(c *ModelConfig, v string) { c.ModelName = v })
	WithBaseURL     = MakeOption(func(c *ModelConfig, v string) { c.BaseURL = v })
	WithAPIKey      = MakeOption(func(c *ModelConfig, v string) { c.APIKey = v })
	WithByAzure     = MakeOption(func(c *ModelConfig, v string) { c.ByAzure = v })
	WithByAzureBool = MakeOption(func(c *ModelConfig, v bool) {
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

func RegisterModelConstructor(modelDesc *ModelDescriber, modelConstructor ModelConstructor) error {
	return modelConstructorRegistry.Register(modelDesc, modelConstructor)
}

// 注册模型构造函数；无特定配置时不要传参
func RegisterModelConstructFunc(modelType ModelType, modelConstructFunc ModelConstructFunc, bindValues ...interface{}) error {
	modelDesc := NewModelDescriber(modelType)
	modelConstructor := NewComponentConstructor[*ModelDescriber, *ModelConfig, model.ToolCallingChatModel](modelDesc, modelConstructFunc)
	return modelConstructorRegistry.Register(modelDesc, modelConstructor, bindValues...)
}

func LookupModelDescriber(value interface{}) (*ModelDescriber, error) {
	return modelConstructorRegistry.LookupDesc(value)
}

func NewChatModel(ctx context.Context, modelOptions ...ModelOption) (model.ToolCallingChatModel, error) {
	modelConfig := NewModelConfig(modelOptions...)
	modelConstructor, err := GetModelConstructor(NewModelDescriber(modelConfig.ModelType))
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

const (
	defaultLocalModelType = "ollama"
	defaultLocalModelName = "qwen2:7b"
	defaultLocalBaseURL   = "http://localhost:11434"
)

// 创建本地模型
func NewLocalChatModel(ctx context.Context, modelOptions ...ModelOption) (model.ToolCallingChatModel, error) {
	// 默认值（均可被用户选项覆盖） + 用户选项
	combinedOptions := make([]ModelOption, 0, len(modelOptions)+3)
	combinedOptions = append(combinedOptions, WithModelType(defaultLocalModelType), WithModelName(defaultLocalModelName), WithBaseURL(defaultLocalBaseURL))
	combinedOptions = append(combinedOptions, modelOptions...)
	return NewChatModel(ctx, combinedOptions...)
}
