// 任务管理
package plantask

import (
	"context"

	"github.com/cloudwego/eino/adk"
	ptmiddleware "github.com/cloudwego/eino/adk/middlewares/plantask"
	"github.com/leitiannet/einolib"
)

type PlanTaskMiddlewareConfig struct {
	ptmiddleware.Config // 内嵌结构体
}

func NewPlanTaskMiddlewareConfig(planTaskMiddlewareOptions ...PlanTaskMiddlewareOption) *PlanTaskMiddlewareConfig {
	planTaskMiddlewareConfig := &PlanTaskMiddlewareConfig{}
	einolib.ApplyOptions(planTaskMiddlewareConfig, planTaskMiddlewareOptions)
	return planTaskMiddlewareConfig
}

type PlanTaskMiddlewareOption func(*PlanTaskMiddlewareConfig)

var (
	WithBackend = einolib.MakeOption(func(c *PlanTaskMiddlewareConfig, v ptmiddleware.Backend) { c.Backend = v })
	WithBaseDir = einolib.MakeOption(func(c *PlanTaskMiddlewareConfig, v string) { c.BaseDir = v })
)

func NewPlanTaskMiddleware(ctx context.Context, planTaskMiddlewareConfig *PlanTaskMiddlewareConfig) (adk.ChatModelAgentMiddleware, error) {
	return ptmiddleware.New(ctx, &planTaskMiddlewareConfig.Config)
}

func createPlanTaskMiddleware(ctx context.Context, middlewareConfig *einolib.MiddlewareConfig, specificConfig interface{}) (adk.ChatModelAgentMiddleware, error) {
	planTaskMiddlewareConfig, err := einolib.ParseSpecificConfig(specificConfig, func() *PlanTaskMiddlewareConfig { return NewPlanTaskMiddlewareConfig() })
	if err != nil {
		return nil, err
	}
	return NewPlanTaskMiddleware(ctx, planTaskMiddlewareConfig)
}

func init() {
	if err := einolib.RegisterMiddlewareConstructFunc(einolib.MiddlewareTypePlanTask, einolib.MiddlewareNameGeneral, createPlanTaskMiddleware, (*PlanTaskMiddlewareConfig)(nil)); err != nil {
		einolib.GetLogger().Errorf("register middleware %s failed: %v", einolib.MiddlewareTypePlanTask, err)
	}
}
