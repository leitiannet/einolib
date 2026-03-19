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

func init() {
	_ = einolib.RegisterModelConstructFunc(ModelTypeOllama, func(ctx context.Context, modelConfig *einolib.ModelConfig, specificConfig interface{}) (model.ToolCallingChatModel, error) {
		return ollama.NewChatModel(ctx, &ollama.ChatModelConfig{
			Model:   modelConfig.ModelName,
			BaseURL: modelConfig.BaseURL,
		})
	})
}
