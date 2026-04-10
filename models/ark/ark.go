package ark

import (
	"context"

	"github.com/cloudwego/eino-ext/components/model/ark"
	"github.com/cloudwego/eino/components/model"
	"github.com/leitiannet/einolib"
)

const (
	ModelTypeARK einolib.ModelType = "ark"
)

func NewARKChatModel(ctx context.Context, modelConfig *einolib.ModelConfig, specificConfig interface{}) (model.ToolCallingChatModel, error) {
	return ark.NewChatModel(ctx, &ark.ChatModelConfig{
		Model:   modelConfig.ModelName,
		BaseURL: modelConfig.BaseURL,
		APIKey:  modelConfig.APIKey,
	})
}

func init() {
	if err := einolib.RegisterModelConstructFunc(ModelTypeARK, NewARKChatModel); err != nil {
		einolib.GetLogger().Errorf("register model %s failed: %v", ModelTypeARK, err)
	}
}
