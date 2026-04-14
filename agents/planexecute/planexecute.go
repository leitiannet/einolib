// Plan-Execute 预构建智能体，组合规划器、执行器与重规划器
package planexecute

import (
	"context"
	"fmt"

	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/adk/prebuilt/planexecute"
	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/schema"
	"github.com/leitiannet/einolib"
)

const (
	AgentTypePlanExecute einolib.AgentType = "planexecute"
)

type PlanExecuteAgentConfig struct {
	planexecute.Config                              // 内嵌planexecute.Config
	PlannerConfig      *planexecute.PlannerConfig   // 规划器配置
	ExecutorConfig     *planexecute.ExecutorConfig  // 执行器配置
	ReplannerConfig    *planexecute.ReplannerConfig // 重规划器配置
}

func NewPlanExecuteAgentConfig(planExecuteAgentOptions ...PlanExecuteAgentOption) *PlanExecuteAgentConfig {
	planExecuteAgentConfig := &PlanExecuteAgentConfig{}
	planExecuteAgentConfig.MaxIterations = 10 // 默认最大迭代次数为10
	einolib.ApplyOptions(planExecuteAgentConfig, planExecuteAgentOptions)
	return planExecuteAgentConfig
}

type PlanExecuteAgentOption func(planExecuteAgentConfig *PlanExecuteAgentConfig)

func initPlannerConfig(c *PlanExecuteAgentConfig) {
	if c.PlannerConfig == nil {
		c.PlannerConfig = &planexecute.PlannerConfig{}
	}
}

func initExecutorConfig(c *PlanExecuteAgentConfig) {
	if c.ExecutorConfig == nil {
		c.ExecutorConfig = &planexecute.ExecutorConfig{}
	}
}

func initReplannerConfig(c *PlanExecuteAgentConfig) {
	if c.ReplannerConfig == nil {
		c.ReplannerConfig = &planexecute.ReplannerConfig{}
	}
}

var (
	// 顶层选项
	WithPlanner         = einolib.MakeOption(func(c *PlanExecuteAgentConfig, v adk.Agent) { c.Planner = v })
	WithExecutor        = einolib.MakeOption(func(c *PlanExecuteAgentConfig, v adk.Agent) { c.Executor = v })
	WithReplanner       = einolib.MakeOption(func(c *PlanExecuteAgentConfig, v adk.Agent) { c.Replanner = v })
	WithPlannerConfig   = einolib.MakeOption(func(c *PlanExecuteAgentConfig, v *planexecute.PlannerConfig) { c.PlannerConfig = v })
	WithExecutorConfig  = einolib.MakeOption(func(c *PlanExecuteAgentConfig, v *planexecute.ExecutorConfig) { c.ExecutorConfig = v })
	WithReplannerConfig = einolib.MakeOption(func(c *PlanExecuteAgentConfig, v *planexecute.ReplannerConfig) { c.ReplannerConfig = v })
	WithMaxIterations   = einolib.MakeOption(func(c *PlanExecuteAgentConfig, v int) { c.MaxIterations = v })

	// Planner 子配置选项（懒初始化）
	WithPlannerChatModelWithFormattedOutput = einolib.MakeOption(func(c *PlanExecuteAgentConfig, v model.BaseChatModel) {
		initPlannerConfig(c)
		c.PlannerConfig.ChatModelWithFormattedOutput = v
	})
	WithPlannerToolCallingChatModel = einolib.MakeOption(func(c *PlanExecuteAgentConfig, v model.ToolCallingChatModel) {
		initPlannerConfig(c)
		c.PlannerConfig.ToolCallingChatModel = v
	})
	WithPlannerToolInfo = einolib.MakeOption(func(c *PlanExecuteAgentConfig, v *schema.ToolInfo) {
		initPlannerConfig(c)
		c.PlannerConfig.ToolInfo = v
	})
	WithPlannerGenInputFn = einolib.MakeOption(func(c *PlanExecuteAgentConfig, v planexecute.GenPlannerModelInputFn) {
		initPlannerConfig(c)
		c.PlannerConfig.GenInputFn = v
	})
	WithPlannerNewPlan = einolib.MakeOption(func(c *PlanExecuteAgentConfig, v planexecute.NewPlan) {
		initPlannerConfig(c)
		c.PlannerConfig.NewPlan = v
	})

	// Executor 子配置选项（懒初始化）
	WithExecutorModel = einolib.MakeOption(func(c *PlanExecuteAgentConfig, v model.ToolCallingChatModel) {
		initExecutorConfig(c)
		c.ExecutorConfig.Model = v
	})
	WithExecutorToolsConfig = einolib.MakeOption(func(c *PlanExecuteAgentConfig, v adk.ToolsConfig) {
		initExecutorConfig(c)
		c.ExecutorConfig.ToolsConfig = v
	})
	WithExecutorMaxIterations = einolib.MakeOption(func(c *PlanExecuteAgentConfig, v int) {
		initExecutorConfig(c)
		c.ExecutorConfig.MaxIterations = v
	})
	WithExecutorGenInputFn = einolib.MakeOption(func(c *PlanExecuteAgentConfig, v planexecute.GenModelInputFn) {
		initExecutorConfig(c)
		c.ExecutorConfig.GenInputFn = v
	})

	// Replanner 子配置选项（懒初始化）
	WithReplannerChatModel = einolib.MakeOption(func(c *PlanExecuteAgentConfig, v model.ToolCallingChatModel) {
		initReplannerConfig(c)
		c.ReplannerConfig.ChatModel = v
	})
	WithReplannerPlanTool = einolib.MakeOption(func(c *PlanExecuteAgentConfig, v *schema.ToolInfo) {
		initReplannerConfig(c)
		c.ReplannerConfig.PlanTool = v
	})
	WithReplannerRespondTool = einolib.MakeOption(func(c *PlanExecuteAgentConfig, v *schema.ToolInfo) {
		initReplannerConfig(c)
		c.ReplannerConfig.RespondTool = v
	})
	WithReplannerGenInputFn = einolib.MakeOption(func(c *PlanExecuteAgentConfig, v planexecute.GenModelInputFn) {
		initReplannerConfig(c)
		c.ReplannerConfig.GenInputFn = v
	})
	WithReplannerNewPlan = einolib.MakeOption(func(c *PlanExecuteAgentConfig, v planexecute.NewPlan) {
		initReplannerConfig(c)
		c.ReplannerConfig.NewPlan = v
	})
)

func newPlanner(ctx context.Context, plannerConfig *planexecute.PlannerConfig) (adk.Agent, error) {
	c := *plannerConfig
	return planexecute.NewPlanner(ctx, &c)
}

func newExecutor(ctx context.Context, executorConfig *planexecute.ExecutorConfig) (adk.Agent, error) {
	c := *executorConfig
	return planexecute.NewExecutor(ctx, &c)
}

func newReplanner(ctx context.Context, replannerConfig *planexecute.ReplannerConfig) (adk.Agent, error) {
	c := *replannerConfig
	return planexecute.NewReplanner(ctx, &c)
}

func NewPlanExecuteAgent(ctx context.Context, agentConfig *einolib.AgentConfig, planExecuteAgentConfig *PlanExecuteAgentConfig) (adk.Agent, error) {
	var err error
	if planExecuteAgentConfig.Planner == nil && planExecuteAgentConfig.PlannerConfig != nil {
		planExecuteAgentConfig.Planner, err = newPlanner(ctx, planExecuteAgentConfig.PlannerConfig)
		if err != nil {
			return nil, fmt.Errorf("planexecute agent: new planner: %v", err)
		}
	}
	if planExecuteAgentConfig.Executor == nil && planExecuteAgentConfig.ExecutorConfig != nil {
		planExecuteAgentConfig.Executor, err = newExecutor(ctx, planExecuteAgentConfig.ExecutorConfig)
		if err != nil {
			return nil, fmt.Errorf("planexecute agent: new executor: %v", err)
		}
	}
	if planExecuteAgentConfig.Replanner == nil && planExecuteAgentConfig.ReplannerConfig != nil {
		planExecuteAgentConfig.Replanner, err = newReplanner(ctx, planExecuteAgentConfig.ReplannerConfig)
		if err != nil {
			return nil, fmt.Errorf("planexecute agent: new replanner: %v", err)
		}
	}
	if planExecuteAgentConfig.Planner == nil {
		return nil, fmt.Errorf("planexecute agent: planner is required")
	}
	if planExecuteAgentConfig.Executor == nil {
		return nil, fmt.Errorf("planexecute agent: executor is required")
	}
	if planExecuteAgentConfig.Replanner == nil {
		return nil, fmt.Errorf("planexecute agent: replanner is required")
	}
	return planexecute.New(ctx, &planExecuteAgentConfig.Config)
}

func createPlanExecuteAgent(ctx context.Context, agentConfig *einolib.AgentConfig, specificConfig interface{}) (adk.Agent, error) {
	planExecuteAgentConfig, err := einolib.ParseSpecificConfig(specificConfig, func() *PlanExecuteAgentConfig { return NewPlanExecuteAgentConfig() })
	if err != nil {
		return nil, err
	}
	return NewPlanExecuteAgent(ctx, agentConfig, planExecuteAgentConfig)
}

func init() {
	if err := einolib.RegisterAgentConstructFunc(AgentTypePlanExecute, einolib.GeneralAgentName, createPlanExecuteAgent, (*PlanExecuteAgentConfig)(nil)); err != nil {
		einolib.GetLogger().Errorf("register agent %s failed: %v", AgentTypePlanExecute, err)
	}
}
