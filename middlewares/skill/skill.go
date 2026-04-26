// 技能中间件：加载并执行SKILL.md定义的技能
package skill

import (
	"context"

	"github.com/cloudwego/eino/adk"
	skmiddleware "github.com/cloudwego/eino/adk/middlewares/skill"
	"github.com/cloudwego/eino/schema"
	"github.com/leitiannet/einolib"
)

type SkillMiddlewareConfig struct {
	skmiddleware.Config // 内嵌结构体
}

func NewSkillMiddlewareConfig(skillMiddlewareOptions ...SkillMiddlewareOption) *SkillMiddlewareConfig {
	skillMiddlewareConfig := &SkillMiddlewareConfig{}
	einolib.ApplyOptions(skillMiddlewareConfig, skillMiddlewareOptions)
	return skillMiddlewareConfig
}

type SkillMiddlewareOption func(*SkillMiddlewareConfig)

var (
	WithBackend               = einolib.MakeOption(func(c *SkillMiddlewareConfig, v skmiddleware.Backend) { c.Backend = v })
	WithSkillToolName         = einolib.MakeOption(func(c *SkillMiddlewareConfig, v string) { c.SkillToolName = &v })
	WithAgentHub              = einolib.MakeOption(func(c *SkillMiddlewareConfig, v skmiddleware.AgentHub) { c.AgentHub = v })
	WithModelHub              = einolib.MakeOption(func(c *SkillMiddlewareConfig, v skmiddleware.ModelHub) { c.ModelHub = v })
	WithCustomSystemPrompt    = einolib.MakeOption(func(c *SkillMiddlewareConfig, v skmiddleware.SystemPromptFunc) { c.CustomSystemPrompt = v })
	WithCustomToolDescription = einolib.MakeOption(func(c *SkillMiddlewareConfig, v skmiddleware.ToolDescriptionFunc) {
		c.CustomToolDescription = v
	})
	WithCustomToolParams = einolib.MakeOption(func(c *SkillMiddlewareConfig, v func(context.Context, map[string]*schema.ParameterInfo) (map[string]*schema.ParameterInfo, error)) {
		c.CustomToolParams = v
	})
	WithBuildContent = einolib.MakeOption(func(c *SkillMiddlewareConfig, v func(context.Context, skmiddleware.Skill, string) (string, error)) {
		c.BuildContent = v
	})
	WithBuildForkMessages = einolib.MakeOption(func(c *SkillMiddlewareConfig, v func(context.Context, skmiddleware.SubAgentInput) ([]adk.Message, error)) {
		c.BuildForkMessages = v
	})
	WithFormatForkResult = einolib.MakeOption(func(c *SkillMiddlewareConfig, v func(context.Context, skmiddleware.SubAgentOutput) (string, error)) {
		c.FormatForkResult = v
	})
)

func NewSkillMiddleware(ctx context.Context, skillMiddlewareConfig *SkillMiddlewareConfig) (adk.ChatModelAgentMiddleware, error) {
	return skmiddleware.NewMiddleware(ctx, &skillMiddlewareConfig.Config)
}

func createSkillMiddleware(ctx context.Context, middlewareConfig *einolib.MiddlewareConfig, specificConfig interface{}) (adk.ChatModelAgentMiddleware, error) {
	skillMiddlewareConfig, err := einolib.ParseSpecificConfig(specificConfig, func() *SkillMiddlewareConfig { return NewSkillMiddlewareConfig() })
	if err != nil {
		return nil, err
	}
	return NewSkillMiddleware(ctx, skillMiddlewareConfig)
}

func init() {
	if err := einolib.RegisterMiddlewareConstructFunc(einolib.MiddlewareTypeSkill, einolib.MiddlewareNameGeneral, createSkillMiddleware, (*SkillMiddlewareConfig)(nil)); err != nil {
		einolib.GetLogger().Errorf("register middleware %s failed: %v", einolib.MiddlewareTypeSkill, err)
	}
}
