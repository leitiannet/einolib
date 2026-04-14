// Ollama 本地部署的聊天模型
package ollama

import (
	"context"

	"github.com/cloudwego/eino-ext/components/model/ollama"
	"github.com/cloudwego/eino/components/model"
	"github.com/leitiannet/einolib"
)

const (
	ModelTypeOllama einolib.ModelType = "ollama"
)

type OllamaModelConfig struct{}

func NewOllamaModelConfig(ollamaModelOptions ...OllamaModelOption) *OllamaModelConfig {
	ollamaModelConfig := &OllamaModelConfig{}
	einolib.ApplyOptions(ollamaModelConfig, ollamaModelOptions)
	return ollamaModelConfig
}

type OllamaModelOption func(*OllamaModelConfig)

func NewOllamaChatModel(ctx context.Context, modelConfig *einolib.ModelConfig, ollamaModelConfig *OllamaModelConfig) (model.ToolCallingChatModel, error) {
	return ollama.NewChatModel(ctx, &ollama.ChatModelConfig{
		Model:   modelConfig.ModelName,
		BaseURL: modelConfig.BaseURL,
	})
}

func createOllamaChatModel(ctx context.Context, modelConfig *einolib.ModelConfig, specificConfig interface{}) (model.ToolCallingChatModel, error) {
	ollamaModelConfig, err := einolib.ParseSpecificConfig(specificConfig, func() *OllamaModelConfig { return NewOllamaModelConfig() })
	if err != nil {
		return nil, err
	}
	return NewOllamaChatModel(ctx, modelConfig, ollamaModelConfig)
}

func init() {
	if err := einolib.RegisterModelConstructFunc(ModelTypeOllama, createOllamaChatModel, (*OllamaModelConfig)(nil)); err != nil {
		einolib.GetLogger().Errorf("register model %s failed: %v", ModelTypeOllama, err)
	}
}
