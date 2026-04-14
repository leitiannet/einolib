// 火山引擎方舟（ARK）聊天模型
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

type ARKModelConfig struct{}

func NewARKModelConfig(arkModelOptions ...ARKModelOption) *ARKModelConfig {
	arkModelConfig := &ARKModelConfig{}
	einolib.ApplyOptions(arkModelConfig, arkModelOptions)
	return arkModelConfig
}

type ARKModelOption func(*ARKModelConfig)

func NewARKChatModel(ctx context.Context, modelConfig *einolib.ModelConfig, arkModelConfig *ARKModelConfig) (model.ToolCallingChatModel, error) {
	return ark.NewChatModel(ctx, &ark.ChatModelConfig{
		Model:   modelConfig.ModelName,
		BaseURL: modelConfig.BaseURL,
		APIKey:  modelConfig.APIKey,
	})
}

func createARKChatModel(ctx context.Context, modelConfig *einolib.ModelConfig, specificConfig interface{}) (model.ToolCallingChatModel, error) {
	arkModelConfig, err := einolib.ParseSpecificConfig(specificConfig, func() *ARKModelConfig { return NewARKModelConfig() })
	if err != nil {
		return nil, err
	}
	return NewARKChatModel(ctx, modelConfig, arkModelConfig)
}

func init() {
	if err := einolib.RegisterModelConstructFunc(ModelTypeARK, createARKChatModel, (*ARKModelConfig)(nil)); err != nil {
		einolib.GetLogger().Errorf("register model %s failed: %v", ModelTypeARK, err)
	}
}
