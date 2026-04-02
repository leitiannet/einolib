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

// 模型元信息
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
		return fmt.Errorf("modelType invalid: %s", md.ModelType)
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

// 将环境变量直接绑定到配置结构体，环境变量前缀：EINO_
// EINO_MODEL_TYPE: 模型类型
// EINO_MODEL_NAME: 模型名称
// EINO_BASE_URL: 服务地址
// EINO_API_KEY: API密钥
func (modelConfig *ModelConfig) initFromEnvironment() {
	if modelConfig == nil {
		return
	}
	_ = envconfig.Process("EINO", modelConfig)
	// 不同模型使用不同环境变量
	BindVarFromEnv((*string)(&modelConfig.ModelType), "MODEL_TYPE")
	if modelConfig.ModelType == ModelTypeUnknown {
		return
	}
	prefix := string(modelConfig.ModelType)
	BindVarFromEnv(&modelConfig.ModelName, "MODEL", prefix)
	BindVarFromEnv(&modelConfig.BaseURL, "BASE_URL", prefix)
	BindVarFromEnv(&modelConfig.APIKey, "API_KEY", prefix)
	BindVarFromEnv(&modelConfig.ByAzure, "BY_AZURE", prefix)
}

// 模型选项
type ModelOption func(modelConfig *ModelConfig)

func WithModelComponentConfig(modelType ModelType, value interface{}) ModelOption {
	return func(modelConfig *ModelConfig) {
		if modelConfig != nil {
			modelConfig.SetConfig(NewModelDescriber(modelType), value)
		}
	}
}

func WithModelType(modelType ModelType) ModelOption {
	return func(modelConfig *ModelConfig) {
		if modelConfig != nil {
			modelConfig.ModelType = modelType
		}
	}
}

func WithModelName(modelName string) ModelOption {
	return func(modelConfig *ModelConfig) {
		if modelConfig != nil {
			modelConfig.ModelName = modelName
		}
	}
}

func WithBaseURL(baseURL string) ModelOption {
	return func(modelConfig *ModelConfig) {
		if modelConfig != nil {
			modelConfig.BaseURL = baseURL
		}
	}
}

func WithAPIKey(apiKey string) ModelOption {
	return func(modelConfig *ModelConfig) {
		if modelConfig != nil {
			modelConfig.APIKey = apiKey
		}
	}
}

func WithByAzure(byAzure string) ModelOption {
	return func(modelConfig *ModelConfig) {
		if modelConfig != nil {
			modelConfig.ByAzure = byAzure
		}
	}
}

func WithByAzureBool(byAzure bool) ModelOption {
	return func(modelConfig *ModelConfig) {
		if modelConfig != nil {
			if byAzure {
				modelConfig.ByAzure = "true"
			} else {
				modelConfig.ByAzure = "false"
			}
		}
	}
}

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

func RegisterModelConstructFunc(modelType ModelType, modelConstructFunc ModelConstructFunc) error {
	modelDesc := NewModelDescriber(modelType)
	modelConstructor := NewComponentConstructor[*ModelDescriber, *ModelConfig, model.ToolCallingChatModel](modelDesc, modelConstructFunc)
	return RegisterModelConstructor(modelDesc, modelConstructor)
}

func GetChatModel(ctx context.Context, modelOptions ...ModelOption) (model.ToolCallingChatModel, error) {
	modelConfig := NewModelConfig(modelOptions...)
	modelConstructor, err := GetModelConstructor(NewModelDescriber(modelConfig.ModelType))
	if err != nil {
		return nil, err
	}
	return modelConstructor.Construct(ctx, modelConfig)
}

func MustGetChatModel(ctx context.Context, modelOptions ...ModelOption) model.ToolCallingChatModel {
	chatModel, err := GetChatModel(ctx, modelOptions...)
	if err != nil {
		panic(fmt.Sprintf("MustGetChatModel failed: %v", err))
	}
	if chatModel == nil {
		panic("MustGetChatModel failed: instance is nil")
	}
	return chatModel
}

// 获取本地ollama模型
func GetLocalChatModel(ctx context.Context, modelOptions ...ModelOption) (model.ToolCallingChatModel, error) {
	// WithModelType("ollama")最后应用，保证ModelType始终为ollama，避免被其他ModelOption覆盖
	combinedOptions := append(
		[]ModelOption{WithModelName("qwen2:7b"), WithBaseURL("http://localhost:11434")}, // 默认值（可被覆盖）
		modelOptions...,
	)
	combinedOptions = append(combinedOptions,
		WithModelType("ollama"),
	)
	return GetChatModel(ctx, combinedOptions...)
}
