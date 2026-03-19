package duckduckgosearch

import (
	"context"
	"fmt"
	"time"

	"github.com/cloudwego/eino-ext/components/tool/duckduckgo/v2"
	"github.com/cloudwego/eino/components/tool"
	"github.com/leitiannet/einolib"
)

const (
	DuckDuckGoSearchToolName = "duckduckgosearch"
)

type DuckDuckGoSearchToolConfig struct {
	Config *duckduckgo.Config // 配置参数
}

func NewDuckDuckGoSearchToolConfig(duckDuckGoSearchToolOptions ...DuckDuckGoSearchToolOption) *DuckDuckGoSearchToolConfig {
	duckDuckGoSearchToolConfig := &DuckDuckGoSearchToolConfig{
		Config: &duckduckgo.Config{
			Region:     duckduckgo.RegionWT,
			Timeout:    10 * time.Second,
			MaxResults: 20, // Limit to return 20 results
		},
	}
	einolib.ApplyOptions(duckDuckGoSearchToolConfig, duckDuckGoSearchToolOptions)
	return duckDuckGoSearchToolConfig
}

type DuckDuckGoSearchToolOption func(*DuckDuckGoSearchToolConfig)

func WithConfig(config *duckduckgo.Config) DuckDuckGoSearchToolOption {
	return func(duckDuckGoSearchToolConfig *DuckDuckGoSearchToolConfig) {
		if duckDuckGoSearchToolConfig != nil {
			duckDuckGoSearchToolConfig.Config = config
		}
	}
}

func getDuckDuckGoSearchTool(ctx context.Context, duckDuckGoSearchToolConfig *DuckDuckGoSearchToolConfig) ([]tool.BaseTool, error) {
	if duckDuckGoSearchToolConfig == nil {
		return nil, fmt.Errorf("duckDuckGoSearchToolConfig nil")
	}
	toolInstance, err := duckduckgo.NewTextSearchTool(ctx, duckDuckGoSearchToolConfig.Config)
	if err != nil {
		return nil, err
	}
	return []tool.BaseTool{toolInstance}, nil
}

func init() {
	_ = einolib.RegisterToolConstructFunc(einolib.ToolTypeBuiltin, DuckDuckGoSearchToolName, func(ctx context.Context, toolConfig *einolib.ToolConfig, specificConfig interface{}) ([]tool.BaseTool, error) {
		duckDuckGoSearchToolConfig := NewDuckDuckGoSearchToolConfig()
		if specificConfig != nil {
			duckDuckGoSearchToolConfig, _ = specificConfig.(*DuckDuckGoSearchToolConfig)
			if duckDuckGoSearchToolConfig == nil {
				return nil, fmt.Errorf("duckDuckGoSearchToolConfig is nil")
			}
		}
		return getDuckDuckGoSearchTool(ctx, duckDuckGoSearchToolConfig)
	})
}
