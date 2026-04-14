// 本地执行器，支持文件操作和命令执行
package local

import (
	"context"

	einolocal "github.com/cloudwego/eino-ext/adk/backend/local"
	"github.com/leitiannet/einolib"
)

const (
	OperatorTypeLocal einolib.OperatorType = "local"
)

type LocalOperatorConfig struct {
	einolocal.Config // 内嵌结构体
}

func NewLocalOperatorConfig(localOperatorOptions ...LocalOperatorOption) *LocalOperatorConfig {
	localOperatorConfig := &LocalOperatorConfig{}
	einolib.ApplyOptions(localOperatorConfig, localOperatorOptions)
	return localOperatorConfig
}

type LocalOperatorOption func(*LocalOperatorConfig)

var (
	WithValidateCommand = einolib.MakeOption(func(c *LocalOperatorConfig, v func(string) error) { c.ValidateCommand = v })
)

func NewLocalOperator(ctx context.Context, localOperatorConfig *LocalOperatorConfig) (einolib.Operator, error) {
	return einolocal.NewBackend(ctx, &localOperatorConfig.Config)
}

func createLocalOperator(ctx context.Context, operatorConfig *einolib.OperatorConfig, specificConfig interface{}) (einolib.Operator, error) {
	localOperatorConfig, err := einolib.ParseSpecificConfig(specificConfig, func() *LocalOperatorConfig { return NewLocalOperatorConfig() })
	if err != nil {
		return nil, err
	}
	return NewLocalOperator(ctx, localOperatorConfig)
}

func init() {
	if err := einolib.RegisterOperatorConstructFunc(OperatorTypeLocal, einolib.GeneralOperatorName, createLocalOperator, (*LocalOperatorConfig)(nil)); err != nil {
		einolib.GetLogger().Errorf("register operator %s failed: %v", OperatorTypeLocal, err)
	}
}
