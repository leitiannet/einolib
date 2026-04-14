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
	MiddlewareTypeFileSystem einolib.MiddlewareType = "filesystem"
)

type FileSystemMiddlewareConfig struct {
	fsmiddleware.MiddlewareConfig // 内嵌结构体
}

func NewFileSystemMiddlewareConfig(fileSystemMiddlewareOptions ...FileSystemMiddlewareOption) *FileSystemMiddlewareConfig {
	fileSystemMiddlewareConfig := &FileSystemMiddlewareConfig{}
	einolib.ApplyOptions(fileSystemMiddlewareConfig, fileSystemMiddlewareOptions)
	return fileSystemMiddlewareConfig
}

type FileSystemMiddlewareOption func(*FileSystemMiddlewareConfig)

var (
	WithBackend             = einolib.MakeOption(func(c *FileSystemMiddlewareConfig, v fsbackend.Backend) { c.Backend = v })
	WithShell               = einolib.MakeOption(func(c *FileSystemMiddlewareConfig, v fsbackend.Shell) { c.Shell = v })
	WithStreamingShell      = einolib.MakeOption(func(c *FileSystemMiddlewareConfig, v fsbackend.StreamingShell) { c.StreamingShell = v })
	WithLsToolConfig        = einolib.MakeOption(func(c *FileSystemMiddlewareConfig, v *fsmiddleware.ToolConfig) { c.LsToolConfig = v })
	WithReadFileToolConfig  = einolib.MakeOption(func(c *FileSystemMiddlewareConfig, v *fsmiddleware.ToolConfig) { c.ReadFileToolConfig = v })
	WithWriteFileToolConfig = einolib.MakeOption(func(c *FileSystemMiddlewareConfig, v *fsmiddleware.ToolConfig) { c.WriteFileToolConfig = v })
	WithEditFileToolConfig  = einolib.MakeOption(func(c *FileSystemMiddlewareConfig, v *fsmiddleware.ToolConfig) { c.EditFileToolConfig = v })
	WithGlobToolConfig      = einolib.MakeOption(func(c *FileSystemMiddlewareConfig, v *fsmiddleware.ToolConfig) { c.GlobToolConfig = v })
	WithGrepToolConfig      = einolib.MakeOption(func(c *FileSystemMiddlewareConfig, v *fsmiddleware.ToolConfig) { c.GrepToolConfig = v })
	WithCustomSystemPrompt  = einolib.MakeOption(func(c *FileSystemMiddlewareConfig, v string) { c.CustomSystemPrompt = &v })
)

func NewFileSystemMiddleware(ctx context.Context, fileSystemMiddlewareConfig *FileSystemMiddlewareConfig) (adk.ChatModelAgentMiddleware, error) {
	return fsmiddleware.New(ctx, &fileSystemMiddlewareConfig.MiddlewareConfig)
}

func createFileSystemMiddleware(ctx context.Context, middlewareConfig *einolib.MiddlewareConfig, specificConfig interface{}) (adk.ChatModelAgentMiddleware, error) {
	fileSystemMiddlewareConfig, err := einolib.ParseSpecificConfig(specificConfig, func() *FileSystemMiddlewareConfig { return NewFileSystemMiddlewareConfig() })
	if err != nil {
		return nil, err
	}
	return NewFileSystemMiddleware(ctx, fileSystemMiddlewareConfig)
}

func init() {
	if err := einolib.RegisterMiddlewareConstructFunc(MiddlewareTypeFileSystem, einolib.GeneralMiddlewareName, createFileSystemMiddleware, (*FileSystemMiddlewareConfig)(nil)); err != nil {
		einolib.GetLogger().Errorf("register middleware %s failed: %v", MiddlewareTypeFileSystem, err)
	}
}
