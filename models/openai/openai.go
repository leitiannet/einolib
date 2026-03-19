package openai

import (
	"context"

	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino/components/model"
	"github.com/leitiannet/einolib"
)

const (
	ModelTypeOpenAI einolib.ModelType = "openai"
)

func init() {
	_ = einolib.RegisterModelConstructFunc(ModelTypeOpenAI, func(ctx context.Context, modelConfig *einolib.ModelConfig, specificConfig interface{}) (model.ToolCallingChatModel, error) {
		return openai.NewChatModel(ctx, &openai.ChatModelConfig{
			Model:   modelConfig.ModelName,
			BaseURL: modelConfig.BaseURL,
			APIKey:  modelConfig.APIKey,
			ByAzure: func() bool {
				return modelConfig.ByAzure == "true"
			}(),
		})
	})
}
