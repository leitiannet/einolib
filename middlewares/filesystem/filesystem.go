// 注入文件操作工具和命令执行工具
package filesystem

import (
	"context"

	"github.com/cloudwego/eino/adk"
	fsbackend "github.com/cloudwego/eino/adk/filesystem"
	fsmiddleware "github.com/cloudwego/eino/adk/middlewares/filesystem"
	"github.com/leitiannet/einolib"
)

const (
	MiddlewareTypeFilesystem einolib.MiddlewareType = "filesystem"
)

type FilesystemMiddlewareConfig struct {
	fsmiddleware.MiddlewareConfig // 内嵌结构体
}

func NewFilesystemMiddlewareConfig(filesystemMiddlewareOptions ...FilesystemMiddlewareOption) *FilesystemMiddlewareConfig {
	config := &FilesystemMiddlewareConfig{}
	einolib.ApplyOptions(config, filesystemMiddlewareOptions)
	return config
}

type FilesystemMiddlewareOption func(*FilesystemMiddlewareConfig)

var (
	WithBackend             = einolib.MakeOption(func(c *FilesystemMiddlewareConfig, v fsbackend.Backend) { c.Backend = v })
	WithShell               = einolib.MakeOption(func(c *FilesystemMiddlewareConfig, v fsbackend.Shell) { c.Shell = v })
	WithStreamingShell      = einolib.MakeOption(func(c *FilesystemMiddlewareConfig, v fsbackend.StreamingShell) { c.StreamingShell = v })
	WithLsToolConfig        = einolib.MakeOption(func(c *FilesystemMiddlewareConfig, v *fsmiddleware.ToolConfig) { c.LsToolConfig = v })
	WithReadFileToolConfig  = einolib.MakeOption(func(c *FilesystemMiddlewareConfig, v *fsmiddleware.ToolConfig) { c.ReadFileToolConfig = v })
	WithWriteFileToolConfig = einolib.MakeOption(func(c *FilesystemMiddlewareConfig, v *fsmiddleware.ToolConfig) { c.WriteFileToolConfig = v })
	WithEditFileToolConfig  = einolib.MakeOption(func(c *FilesystemMiddlewareConfig, v *fsmiddleware.ToolConfig) { c.EditFileToolConfig = v })
	WithGlobToolConfig      = einolib.MakeOption(func(c *FilesystemMiddlewareConfig, v *fsmiddleware.ToolConfig) { c.GlobToolConfig = v })
	WithGrepToolConfig      = einolib.MakeOption(func(c *FilesystemMiddlewareConfig, v *fsmiddleware.ToolConfig) { c.GrepToolConfig = v })
	WithCustomSystemPrompt  = einolib.MakeOption(func(c *FilesystemMiddlewareConfig, v string) { c.CustomSystemPrompt = &v })
)

func NewFilesystemMiddleware(ctx context.Context, config *FilesystemMiddlewareConfig) (adk.ChatModelAgentMiddleware, error) {
	return fsmiddleware.New(ctx, &config.MiddlewareConfig)
}

func createMiddleware(ctx context.Context, middlewareConfig *einolib.MiddlewareConfig, specificConfig interface{}) (adk.ChatModelAgentMiddleware, error) {
	filesystemMiddlewareConfig, err := einolib.ParseSpecificConfig(specificConfig, func() *FilesystemMiddlewareConfig { return NewFilesystemMiddlewareConfig() })
	if err != nil {
		return nil, err
	}
	return NewFilesystemMiddleware(ctx, filesystemMiddlewareConfig)
}

func init() {
	if err := einolib.RegisterMiddlewareConstructFunc(MiddlewareTypeFilesystem, einolib.GeneralMiddlewareName, createMiddleware, (*FilesystemMiddlewareConfig)(nil)); err != nil {
		einolib.GetLogger().Errorf("register middleware %s failed: %v", MiddlewareTypeFilesystem, err)
	}
}
