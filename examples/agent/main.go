package main

import (
	"context"
	"fmt"

	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
	"github.com/leitiannet/einolib"
	"github.com/leitiannet/einolib/agents/chatmodel"
	mwlog "github.com/leitiannet/einolib/middlewares/log"
	_ "github.com/leitiannet/einolib/models"
	_ "github.com/leitiannet/einolib/tools"
	"github.com/leitiannet/einolib/tools/custom/booksearch"
)

func main() {
	ctx := context.Background()
	// 获取工具
	tools, _, err := einolib.GetTool(ctx,
		einolib.WithToolType(einolib.ToolTypeCustom),
		einolib.WithToolName(booksearch.SearchBookToolName),
	)
	if err != nil {
		panic(err)
	}
	if len(tools) == 0 {
		panic("empty tool list")
	}
	// 创建模型
	chatModel, err := einolib.NewLocalChatModel(ctx)
	if err != nil {
		panic(err)
	}
	// 创建智能体
	agent, err := einolib.NewAgent(ctx,
		einolib.WithAgentType(chatmodel.AgentTypeChatModel),
		einolib.WithAgentName("BookRecommender"),
		einolib.WithAgentComponentConfig(chatmodel.NewChatModelAgentConfig(
			chatmodel.WithDescription("An agent that can recommend books"),
			chatmodel.WithInstruction(`You are an expert book recommender. Based on the user's request, use the "search_book" tool to find relevant books. Finally, present the results to the user.`),
			chatmodel.WithModel(chatModel),
			chatmodel.WithToolsConfig(adk.ToolsConfig{
				ToolsNodeConfig: compose.ToolsNodeConfig{Tools: tools},
			}),
			chatmodel.WithMiddlewares(mwlog.AgentMiddleware()),
			chatmodel.WithHandlers(mwlog.NewChatModelAgentMiddleware()),
		)),
	)
	if err != nil {
		panic(err)
	}
	// 创建运行器
	runner, err := einolib.NewRunner(ctx, einolib.WithAgent(agent))
	if err != nil {
		panic(err)
	}
	// 运行
	iter := runner.Query(ctx, "recommend a fiction book to me")
	for {
		event, ok := iter.Next()
		if !ok {
			break
		}
		if event == nil {
			continue
		}
		if event.Err != nil {
			panic(event.Err)
		}
		if event.Output == nil || event.Output.MessageOutput == nil {
			continue
		}
		msg, err := event.Output.MessageOutput.GetMessage()
		if err != nil {
			panic(err)
		}
		if msg == nil {
			continue
		}
		// 仅打印助手最终文本，避免把 Tool 等中间消息混进输出
		if event.Output.MessageOutput.Role != schema.Assistant {
			continue
		}
		fmt.Printf("\nmessage:\n%v\n======\n", msg)
	}
}
