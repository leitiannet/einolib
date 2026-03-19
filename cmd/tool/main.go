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
	mustPrintToolInfos(einolib.ToolTypeCustom, todo.CoreTodoToolName)
	mustPrintToolInfos(einolib.ToolTypeBuiltin, duckduckgosearch.DuckDuckGoSearchToolName)
	mcpToolConfig := einolibmcp.NewMCPToolConfig(einolibmcp.WithBaseURL("http://localhost:3000/sse"))
	mustPrintToolInfos(einolib.ToolTypeMCP, "mcp_weather_tool",
		einolib.WithToolComponentConfig(einolib.ToolTypeMCP, einolib.GeneralToolName, mcpToolConfig))
}

func mustPrintToolInfos(toolType einolib.ToolType, toolName string, opts ...einolib.ToolOption) {
	_, infos, err := einolib.GetTool(context.TODO(), toolType, toolName, opts...)
	if err != nil {
		panic(err)
	}
	einolib.PrintJSON(infos, einolib.NewPrintJSONOptions(fmt.Sprintf("%s tool infos:", toolName), true))
}
