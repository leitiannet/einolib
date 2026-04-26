// OpenAI 兼容的聊天模型（含 Azure OpenAI）
package openai

import (
	"context"

	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino/components/model"
	"github.com/leitiannet/einolib"
)

type OpenAIModelConfig struct{}

func NewOpenAIModelConfig(openAIModelOptions ...OpenAIModelOption) *OpenAIModelConfig {
	openAIModelConfig := &OpenAIModelConfig{}
	einolib.ApplyOptions(openAIModelConfig, openAIModelOptions)
	return openAIModelConfig
}

type OpenAIModelOption func(*OpenAIModelConfig)

func NewOpenAIChatModel(ctx context.Context, modelConfig *einolib.ModelConfig, openAIModelConfig *OpenAIModelConfig) (model.ToolCallingChatModel, error) {
	return openai.NewChatModel(ctx, &openai.ChatModelConfig{
		Model:   modelConfig.ID,
		BaseURL: modelConfig.BaseURL,
		APIKey:  modelConfig.APIKey,
		ByAzure: modelConfig.ByAzure == "true",
	})
}

func createOpenAIChatModel(ctx context.Context, modelConfig *einolib.ModelConfig, specificConfig interface{}) (model.ToolCallingChatModel, error) {
	openAIModelConfig, err := einolib.ParseSpecificConfig(specificConfig, func() *OpenAIModelConfig { return NewOpenAIModelConfig() })
	if err != nil {
		return nil, err
	}
	return NewOpenAIChatModel(ctx, modelConfig, openAIModelConfig)
}

func init() {
	if err := einolib.RegisterModelConstructFunc(einolib.ModelTypeOpenAI, createOpenAIChatModel, (*OpenAIModelConfig)(nil)); err != nil {
		einolib.GetLogger().Errorf("register model %s failed: %v", einolib.ModelTypeOpenAI, err)
	}
}
