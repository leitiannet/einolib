package duckduckgosearch

import (
	"context"
	"time"

	"github.com/cloudwego/eino-ext/components/tool/duckduckgo/v2"
	"github.com/cloudwego/eino/components/tool"
	"github.com/leitiannet/einolib"
)

const (
	DuckDuckGoSearchToolName = "duckduckgosearch"
)

type DuckDuckGoSearchToolConfig struct {
	duckduckgo.Config // 内嵌duckduckgo.Config
}

func NewDuckDuckGoSearchToolConfig(duckDuckGoSearchToolOptions ...DuckDuckGoSearchToolOption) *DuckDuckGoSearchToolConfig {
	duckDuckGoSearchToolConfig := &DuckDuckGoSearchToolConfig{}
	duckDuckGoSearchToolConfig.Region = duckduckgo.RegionWT
	duckDuckGoSearchToolConfig.Timeout = 10 * time.Second
	duckDuckGoSearchToolConfig.MaxResults = 20
	einolib.ApplyOptions(duckDuckGoSearchToolConfig, duckDuckGoSearchToolOptions)
	return duckDuckGoSearchToolConfig
}

type DuckDuckGoSearchToolOption func(duckDuckGoSearchToolConfig *DuckDuckGoSearchToolConfig)

var (
	WithToolName   = einolib.MakeOption(func(c *DuckDuckGoSearchToolConfig, v string) { c.ToolName = v })
	WithToolDesc   = einolib.MakeOption(func(c *DuckDuckGoSearchToolConfig, v string) { c.ToolDesc = v })
	WithTimeout    = einolib.MakeOption(func(c *DuckDuckGoSearchToolConfig, v time.Duration) { c.Timeout = v })
	WithMaxResults = einolib.MakeOption(func(c *DuckDuckGoSearchToolConfig, v int) { c.MaxResults = v })
	WithRegion     = einolib.MakeOption(func(c *DuckDuckGoSearchToolConfig, v duckduckgo.Region) { c.Region = v })
)

func NewDuckDuckGoSearchTool(ctx context.Context, toolConfig *einolib.ToolConfig, specificConfig interface{}) ([]tool.BaseTool, error) {
	duckDuckGoSearchToolConfig, err := einolib.ParseSpecificConfig(specificConfig, func() *DuckDuckGoSearchToolConfig { return NewDuckDuckGoSearchToolConfig() })
	if err != nil {
		return nil, err
	}
	toolInstance, err := duckduckgo.NewTextSearchTool(ctx, &duckDuckGoSearchToolConfig.Config)
	if err != nil {
		return nil, err
	}
	return []tool.BaseTool{toolInstance}, nil
}

func init() {
	if err := einolib.RegisterToolConstructFunc(einolib.ToolTypeBuiltin, DuckDuckGoSearchToolName, NewDuckDuckGoSearchTool, (*DuckDuckGoSearchToolConfig)(nil)); err != nil {
		einolib.GetLogger().Errorf("register tool %s failed: %v", DuckDuckGoSearchToolName, err)
	}
}
