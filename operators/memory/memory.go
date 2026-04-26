// 内存执行器，支持文件操作，但不支持命令执行
package memory

import (
	"context"

	"github.com/cloudwego/eino/adk/filesystem"
	"github.com/leitiannet/einolib"
)

type MemoryOperatorConfig struct{}

func NewMemoryOperatorConfig(memoryOperatorOptions ...MemoryOperatorOption) *MemoryOperatorConfig {
	memoryOperatorConfig := &MemoryOperatorConfig{}
	einolib.ApplyOptions(memoryOperatorConfig, memoryOperatorOptions)
	return memoryOperatorConfig
}

type MemoryOperatorOption func(*MemoryOperatorConfig)

func NewMemoryOperator(ctx context.Context, memoryOperatorConfig *MemoryOperatorConfig) (einolib.Operator, error) {
	return &MemoryOperator{InMemoryBackend: filesystem.NewInMemoryBackend()}, nil
}

func createMemoryOperator(ctx context.Context, operatorConfig *einolib.OperatorConfig, specificConfig interface{}) (einolib.Operator, error) {
	memoryOperatorConfig, err := einolib.ParseSpecificConfig(specificConfig, func() *MemoryOperatorConfig { return NewMemoryOperatorConfig() })
	if err != nil {
		return nil, err
	}
	return NewMemoryOperator(ctx, memoryOperatorConfig)
}

func init() {
	if err := einolib.RegisterOperatorConstructFunc(einolib.OperatorTypeMemory, einolib.OperatorNameGeneral, createMemoryOperator, (*MemoryOperatorConfig)(nil)); err != nil {
		einolib.GetLogger().Errorf("register operator %s failed: %v", einolib.OperatorTypeMemory, err)
	}
}
