// 文件操作和命令执行中间件，每个工具均可独立配置名称、描述、自定义实现或禁用，工具描述和系统提示词支持中英文切换
// 文件操作：ls、read_file、write_file、edit_file、glob、grep
// 命令执行：execute
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

func withToolDisabled(tc *fsmiddleware.ToolConfig) *fsmiddleware.ToolConfig {
	if tc == nil {
		tc = &fsmiddleware.ToolConfig{}
	}
	tc.Disable = true
	return tc
}

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
	// 与With*ToolConfig中设置Disable: true等价
	WithLsToolDisabled = FileSystemMiddlewareOption(func(c *FileSystemMiddlewareConfig) {
		c.LsToolConfig = withToolDisabled(c.LsToolConfig)
	})
	WithReadFileToolDisabled = FileSystemMiddlewareOption(func(c *FileSystemMiddlewareConfig) {
		c.ReadFileToolConfig = withToolDisabled(c.ReadFileToolConfig)
	})
	WithWriteFileToolDisabled = FileSystemMiddlewareOption(func(c *FileSystemMiddlewareConfig) {
		c.WriteFileToolConfig = withToolDisabled(c.WriteFileToolConfig)
	})
	WithEditFileToolDisabled = FileSystemMiddlewareOption(func(c *FileSystemMiddlewareConfig) {
		c.EditFileToolConfig = withToolDisabled(c.EditFileToolConfig)
	})
	WithGlobToolDisabled = FileSystemMiddlewareOption(func(c *FileSystemMiddlewareConfig) {
		c.GlobToolConfig = withToolDisabled(c.GlobToolConfig)
	})
	WithGrepToolDisabled = FileSystemMiddlewareOption(func(c *FileSystemMiddlewareConfig) {
		c.GrepToolConfig = withToolDisabled(c.GrepToolConfig)
	})
	WithAllFileToolsDisabled = FileSystemMiddlewareOption(func(c *FileSystemMiddlewareConfig) {
		c.LsToolConfig = withToolDisabled(c.LsToolConfig)
		c.ReadFileToolConfig = withToolDisabled(c.ReadFileToolConfig)
		c.WriteFileToolConfig = withToolDisabled(c.WriteFileToolConfig)
		c.EditFileToolConfig = withToolDisabled(c.EditFileToolConfig)
		c.GlobToolConfig = withToolDisabled(c.GlobToolConfig)
		c.GrepToolConfig = withToolDisabled(c.GrepToolConfig)
	})
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
