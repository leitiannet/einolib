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

func NewOpenAIChatModel(ctx context.Context, modelConfig *einolib.ModelConfig, specificConfig interface{}) (model.ToolCallingChatModel, error) {
	return openai.NewChatModel(ctx, &openai.ChatModelConfig{
		Model:   modelConfig.ModelName,
		BaseURL: modelConfig.BaseURL,
		APIKey:  modelConfig.APIKey,
		ByAzure: modelConfig.ByAzure == "true",
	})
}

func init() {
	if err := einolib.RegisterModelConstructFunc(ModelTypeOpenAI, NewOpenAIChatModel); err != nil {
		einolib.GetLogger().Errorf("register model %s failed: %v", ModelTypeOpenAI, err)
	}
}
