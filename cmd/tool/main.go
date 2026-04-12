package main

import (
	"context"
	"fmt"

	"github.com/leitiannet/einolib"
	_ "github.com/leitiannet/einolib/tools"
	"github.com/leitiannet/einolib/tools/builtin/duckduckgosearch"
	"github.com/leitiannet/einolib/tools/custom/todo"
	einolibmcp "github.com/leitiannet/einolib/tools/mcp"
)

func main() {
	ctx := einolib.WithCloser(context.Background())
	defer einolib.CloseCloser(ctx)

	mustPrintToolInfos(ctx, einolib.ToolTypeCustom, todo.CoreTodoToolName)
	mustPrintToolInfos(ctx, einolib.ToolTypeBuiltin, duckduckgosearch.DuckDuckGoSearchToolName)
	mcpToolConfig := einolibmcp.NewMCPToolConfig(einolibmcp.WithBaseURL("http://localhost:3000/sse"))
	mustPrintToolInfos(ctx, einolib.ToolTypeMCP, "mcp_weather_tool",
		einolib.WithToolComponentConfig(mcpToolConfig))
}

func mustPrintToolInfos(ctx context.Context, toolType einolib.ToolType, toolName string, opts ...einolib.ToolOption) {
	opts = append([]einolib.ToolOption{einolib.WithToolType(toolType), einolib.WithToolName(toolName)}, opts...)
	_, infos, err := einolib.GetTool(ctx, opts...)
	if err != nil {
		panic(err)
	}
	einolib.PrintJSON(infos, einolib.NewPrintJSONOptions(fmt.Sprintf("%s tool infos:", toolName), true))
}
