package main

import (
	"context"
	"fmt"

	"github.com/cloudwego/eino/schema"
	"github.com/leitiannet/einolib"
	"github.com/leitiannet/einolib/builder"
)

func main() {
	ctx := context.Background()
	// 创建智能体
	agent, err := builder.NewAgentBuilder(
		einolib.AgentTypeChatModel,
		"FileSystemInspector",
		"inspect directory and read main.go with filesystem middleware",
	).WithInstruction("你需要先查看当前目录下文件列表。如果存在 main.go，请读取并返回其内容。").Build(ctx)
	if err != nil {
		panic(err)
	}
	// 创建运行器
	runner, err := einolib.NewRunner(ctx, einolib.WithAgent(agent))
	if err != nil {
		panic(err)
	}

	iter := runner.Query(ctx, "请执行文件检查任务。")
	for {
		event, ok := iter.Next()
		if !ok {
			break
		}
		if event == nil || event.Err != nil || event.Output == nil || event.Output.MessageOutput == nil {
			if event != nil && event.Err != nil {
				panic(event.Err)
			}
			continue
		}
		if event.Output.MessageOutput.Role != schema.Assistant {
			continue
		}
		msg, err := event.Output.MessageOutput.GetMessage()
		if err != nil {
			panic(err)
		}
		if msg == nil {
			continue
		}
		fmt.Printf("\nmessage:\n%v\n======\n", msg)
	}
}
