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

type MCPToolConfig struct {
	Config           *mcpp.Config             // 配置参数
	BaseURL          string                   // 服务地址
	Options          []transport.ClientOption // 客户端选项
	InitializeParams *mcp.InitializeParams    // 初始化参数
}

func NewMCPToolConfig(mcpToolOptions ...MCPToolOption) *MCPToolConfig {
	mcpToolConfig := &MCPToolConfig{
		BaseURL: defaultMCPBaseURL,
		InitializeParams: &mcp.InitializeParams{
			ProtocolVersion: mcp.LATEST_PROTOCOL_VERSION,
			ClientInfo: mcp.Implementation{
				Name:    defaultMCPClientName,
				Version: defaultMCPClientVersion,
			},
		},
	}
	einolib.ApplyOptions(mcpToolConfig, mcpToolOptions)
	return mcpToolConfig
}

type MCPToolOption func(*MCPToolConfig)

func WithConfig(config *mcpp.Config) MCPToolOption {
	return func(mcpToolConfig *MCPToolConfig) {
		if mcpToolConfig != nil {
			mcpToolConfig.Config = config
		}
	}
}

func WithBaseURL(baseURL string) MCPToolOption {
	return func(mcpToolConfig *MCPToolConfig) {
		if mcpToolConfig != nil {
			mcpToolConfig.BaseURL = baseURL
		}
	}
}

func WithOptions(options []transport.ClientOption) MCPToolOption {
	return func(mcpToolConfig *MCPToolConfig) {
		if mcpToolConfig != nil {
			mcpToolConfig.Options = options
		}
	}
}

func WithInitializeParams(initializeParams *mcp.InitializeParams) MCPToolOption {
	return func(mcpToolConfig *MCPToolConfig) {
		if mcpToolConfig != nil {
			mcpToolConfig.InitializeParams = initializeParams
		}
	}
}

func getMCPClient(ctx context.Context, mcpToolConfig *MCPToolConfig) (client.MCPClient, error) {
	if mcpToolConfig == nil {
		return nil, fmt.Errorf("mcpToolConfig nil")
	}
	cli, err := client.NewSSEMCPClient(mcpToolConfig.BaseURL, mcpToolConfig.Options...)
	if err != nil {
		return nil, err
	}
	err = cli.Start(ctx)
	if err != nil {
		return nil, err
	}
	initRequest := mcp.InitializeRequest{
		Params: *mcpToolConfig.InitializeParams,
	}
	_, err = cli.Initialize(ctx, initRequest)
	if err != nil {
		return nil, err
	}
	return cli, nil
}

func getMCPTool(ctx context.Context, mcpToolConfig *MCPToolConfig) ([]tool.BaseTool, error) {
	if mcpToolConfig == nil {
		return nil, fmt.Errorf("mcpToolConfig nil")
	}
	config := mcpToolConfig.Config
	if config == nil {
		config = &mcpp.Config{}
	}
	if config.Cli == nil {
		cli, err := getMCPClient(ctx, mcpToolConfig)
		if err != nil {
			return nil, err
		}
		config.Cli = cli
	}
	tools, err := mcpp.GetTools(ctx, config)
	if err != nil {
		return nil, err
	}
	return tools, nil
}

func init() {
	_ = einolib.RegisterToolConstructFunc(einolib.ToolTypeMCP, einolib.GeneralToolName, func(ctx context.Context, toolConfig *einolib.ToolConfig, specificConfig interface{}) ([]tool.BaseTool, error) {
		mcpToolConfig := NewMCPToolConfig()
		if specificConfig != nil {
			mcpToolConfig, _ = specificConfig.(*MCPToolConfig)
			if mcpToolConfig == nil {
				return nil, fmt.Errorf("mcpToolConfig is nil")
			}
		}
		return getMCPTool(ctx, mcpToolConfig)
	})
}
