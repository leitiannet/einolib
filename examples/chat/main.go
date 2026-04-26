package main

import (
	"context"

	"github.com/cloudwego/eino/schema"
	"github.com/leitiannet/einolib"
	"github.com/leitiannet/einolib/builder"
	_ "github.com/leitiannet/einolib/models"
)

func main() {
	templates := einolib.NewTemplates(
		einolib.Msg(schema.System, "你是一个{role}。你需要用{style}的语气回答问题。你的目标是帮助程序员保持积极乐观的心态，提供技术建议的同时也要关注他们的心理健康。"),
		einolib.Placeholder("chat_history"),
		einolib.Msg(schema.User, "问题: {question}"),
	)
	chatHistory := einolib.NewMessages(
		einolib.Msg(schema.User, "你好"),
		einolib.Msg(schema.Assistant, "嘿！我是你的程序员鼓励师！记住，每个优秀的程序员都是从 Debug 中成长起来的。有什么我可以帮你的吗？"),
		einolib.Msg(schema.User, "我觉得自己写的代码太烂了"),
		einolib.Msg(schema.Assistant, "每个程序员都经历过这个阶段！重要的是你在不断学习和进步。让我们一起看看代码，我相信通过重构和优化，它会变得更好。记住，Rome wasn't built in a day，代码质量是通过持续改进来提升的。"),
	)
	variables := map[string]any{
		"role":         "程序员鼓励师",
		"style":        "积极、温暖且专业",
		"question":     "我的代码一直报错，感觉好沮丧，该怎么办？",
		"chat_history": chatHistory,
	}
	ctx := context.Background()
	messages, err := einolib.FormatMessages(ctx, schema.FString, templates, variables)
	if err != nil {
		panic(err)
	}
	chatModel, err := builder.NewModelBuilder(builder.DefaultLocalModelType).Build(ctx)
	if err != nil {
		panic(err)
	}
	result, err := chatModel.Generate(ctx, messages)
	if err != nil {
		panic(err)
	}
	einolib.PrintJSON(result, einolib.NewPrintJSONOptions("chatModel generate result:", true))
}
