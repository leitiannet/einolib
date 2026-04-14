// MCP 协议工具，支持 SSE、Stdio、Streamable HTTP 等传输方式
package mcp

import (
	"context"
	"fmt"

	mcpp "github.com/cloudwego/eino-ext/components/tool/mcp"
	"github.com/cloudwego/eino/components/tool"
	"github.com/leitiannet/einolib"
	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/client/transport"
	"github.com/mark3labs/mcp-go/mcp"
)

const (
	defaultMCPBaseURL       = "http://localhost:12345/sse"
	defaultMCPClientName    = "einolib-mcp-client"
	defaultMCPClientVersion = "1.0.0"
)

// MCP传输类型
type MCPTransportType string

const (
	MCPTransportSSE            MCPTransportType = "sse"             // SSE传输（默认）
	MCPTransportStdio          MCPTransportType = "stdio"           // 标准输入输出传输
	MCPTransportStreamableHTTP MCPTransportType = "streamable_http" // Streamable HTTP传输
)

type MCPToolConfig struct {
	mcpp.Config                                          // 内嵌mcpp.Config
	TransportType      MCPTransportType                  // 传输类型（默认SSE）
	BaseURL            string                            // SSE/StreamableHTTP 服务地址
	Options            []transport.ClientOption          // SSE 客户端选项
	StreamableHTTPOpts []transport.StreamableHTTPCOption // StreamableHTTP 客户端选项
	Command            string                            // Stdio 命令路径
	Env                []string                          // Stdio 环境变量
	Args               []string                          // Stdio 命令参数
	InitializeParams   *mcp.InitializeParams             // 初始化参数
}

func NewMCPToolConfig(mcpToolOptions ...MCPToolOption) *MCPToolConfig {
	mcpToolConfig := &MCPToolConfig{}
	mcpToolConfig.TransportType = MCPTransportSSE
	mcpToolConfig.BaseURL = defaultMCPBaseURL
	mcpToolConfig.InitializeParams = &mcp.InitializeParams{
		ProtocolVersion: mcp.LATEST_PROTOCOL_VERSION,
		ClientInfo: mcp.Implementation{
			Name:    defaultMCPClientName,
			Version: defaultMCPClientVersion,
		},
	}
	einolib.ApplyOptions(mcpToolConfig, mcpToolOptions)
	return mcpToolConfig
}

type MCPToolOption func(mcpToolConfig *MCPToolConfig)

var (
	WithCli                   = einolib.MakeOption(func(c *MCPToolConfig, v client.MCPClient) { c.Cli = v })
	WithToolNameList          = einolib.MakeOption(func(c *MCPToolConfig, v []string) { c.ToolNameList = v })
	WithCustomHeaders         = einolib.MakeOption(func(c *MCPToolConfig, v map[string]string) { c.CustomHeaders = v })
	WithMeta                  = einolib.MakeOption(func(c *MCPToolConfig, v *mcp.Meta) { c.Meta = v })
	WithBaseURL               = einolib.MakeOption(func(c *MCPToolConfig, v string) { c.BaseURL = v })
	WithOptions               = einolib.MakeOption(func(c *MCPToolConfig, v []transport.ClientOption) { c.Options = v })
	WithStreamableHTTPOpts    = einolib.MakeOption(func(c *MCPToolConfig, v []transport.StreamableHTTPCOption) { c.StreamableHTTPOpts = v })
	WithInitializeParams      = einolib.MakeOption(func(c *MCPToolConfig, v *mcp.InitializeParams) { c.InitializeParams = v })
	WithTransportType         = einolib.MakeOption(func(c *MCPToolConfig, v MCPTransportType) { c.TransportType = v })
	WithCommand               = einolib.MakeOption(func(c *MCPToolConfig, v string) { c.Command = v })
	WithEnv                   = einolib.MakeOption(func(c *MCPToolConfig, v []string) { c.Env = v })
	WithArgs                  = einolib.MakeOption(func(c *MCPToolConfig, v []string) { c.Args = v })
	WithToolCallResultHandler = einolib.MakeOption(func(c *MCPToolConfig, v func(ctx context.Context, name string, result *mcp.CallToolResult) (*mcp.CallToolResult, error)) {
		c.ToolCallResultHandler = v
	})
)

func newMCPClient(ctx context.Context, mcpToolConfig *MCPToolConfig) (client.MCPClient, error) {
	var (
		cli *client.Client
		err error
	)
	switch mcpToolConfig.TransportType {
	case MCPTransportStdio:
		cli, err = client.NewStdioMCPClient(mcpToolConfig.Command, mcpToolConfig.Env, mcpToolConfig.Args...)
	case MCPTransportStreamableHTTP:
		cli, err = client.NewStreamableHttpClient(mcpToolConfig.BaseURL, mcpToolConfig.StreamableHTTPOpts...)
	default:
		cli, err = client.NewSSEMCPClient(mcpToolConfig.BaseURL, mcpToolConfig.Options...)
	}
	if err != nil {
		return nil, err
	}
	// SSE 和 StreamableHTTP 需要手动启动
	if mcpToolConfig.TransportType != MCPTransportStdio {
		if err = cli.Start(ctx); err != nil {
			_ = cli.Close()
			return nil, err
		}
	}
	initRequest := mcp.InitializeRequest{
		Params: *mcpToolConfig.InitializeParams,
	}
	if _, err = cli.Initialize(ctx, initRequest); err != nil {
		_ = cli.Close()
		return nil, err
	}
	return cli, nil
}

type MCPToolResult struct {
	Tools []tool.BaseTool
	cli   client.MCPClient // 内部创建的客户端，需调用方负责关闭（外部注入的客户端不受影响）
}

func (r *MCPToolResult) Close() error {
	if r.cli != nil {
		return r.cli.Close()
	}
	return nil
}

func newMCPTools(ctx context.Context, mcpToolConfig *MCPToolConfig) (*MCPToolResult, error) {
	result := &MCPToolResult{}
	if mcpToolConfig.Cli == nil {
		cli, err := newMCPClient(ctx, mcpToolConfig)
		if err != nil {
			return nil, err
		}
		mcpToolConfig.Cli = cli
		result.cli = cli
	}
	tools, err := mcpp.GetTools(ctx, &mcpToolConfig.Config)
	if err != nil {
		result.Close()
		return nil, err
	}
	result.Tools = tools
	return result, nil
}

func NewMCPTool(ctx context.Context, toolConfig *einolib.ToolConfig, mcpToolConfig *MCPToolConfig) ([]tool.BaseTool, error) {
	switch mcpToolConfig.TransportType {
	case MCPTransportStdio:
		if mcpToolConfig.Command == "" {
			return nil, fmt.Errorf("mcp tool: command is required for stdio transport")
		}
	default:
		if mcpToolConfig.BaseURL == "" {
			return nil, fmt.Errorf("mcp tool: baseURL is required")
		}
	}
	result, err := newMCPTools(ctx, mcpToolConfig)
	if err != nil {
		return nil, err
	}
	if result.cli != nil {
		einolib.AddCloser(ctx, result)
	}
	return result.Tools, nil
}

func createMCPTool(ctx context.Context, toolConfig *einolib.ToolConfig, specificConfig interface{}) ([]tool.BaseTool, error) {
	mcpToolConfig, err := einolib.ParseSpecificConfig(specificConfig, func() *MCPToolConfig { return NewMCPToolConfig() })
	if err != nil {
		return nil, err
	}
	return NewMCPTool(ctx, toolConfig, mcpToolConfig)
}

func init() {
	if err := einolib.RegisterToolConstructFunc(einolib.ToolTypeMCP, einolib.GeneralToolName, createMCPTool, (*MCPToolConfig)(nil)); err != nil {
		einolib.GetLogger().Errorf("register tool %s failed: %v", einolib.GeneralToolName, err)
	}
}
