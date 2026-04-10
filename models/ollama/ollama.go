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

func NewOllamaChatModel(ctx context.Context, modelConfig *einolib.ModelConfig, specificConfig interface{}) (model.ToolCallingChatModel, error) {
	return ollama.NewChatModel(ctx, &ollama.ChatModelConfig{
		Model:   modelConfig.ModelName,
		BaseURL: modelConfig.BaseURL,
	})
}

func init() {
	if err := einolib.RegisterModelConstructFunc(ModelTypeOllama, NewOllamaChatModel); err != nil {
		einolib.GetLogger().Errorf("register model %s failed: %v", ModelTypeOllama, err)
	}
}
