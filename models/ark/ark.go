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

func init() {
	_ = einolib.RegisterModelConstructFunc(ModelTypeARK, func(ctx context.Context, modelConfig *einolib.ModelConfig, specificConfig interface{}) (model.ToolCallingChatModel, error) {
		return ark.NewChatModel(ctx, &ark.ChatModelConfig{
			Model:   modelConfig.ModelName,
			BaseURL: modelConfig.BaseURL,
			APIKey:  modelConfig.APIKey,
		})
	})
}
