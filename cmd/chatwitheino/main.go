package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/cloudwego/eino/schema"
	"github.com/leitiannet/einolib"
	_ "github.com/leitiannet/einolib/models"
)

// EINO_MODEL_TYPE=ollama EINO_MODEL_NAME=qwen2:7b EINO_BASE_URL=http://localhost:11434 go run cmd/chatwitheino/main.go -- "用一句话解释 Eino 的 Component 设计解决了什么问题？"
func main() {
	// 系统指令
	var instruction string
	flag.StringVar(&instruction, "instruction", "You are a helpful assistant.", "system instruction for the assistant")
	flag.Parse()
	// 查询内容
	query := strings.TrimSpace(strings.Join(flag.Args(), " "))
	if query == "" {
		_, _ = fmt.Fprintln(os.Stderr, "usage: go run main.go -- \"your question\"")
		os.Exit(2)
	}
	messages := []*schema.Message{
		schema.SystemMessage(instruction),
		schema.UserMessage(query),
	}
	//
	ctx := context.Background()
	cm, err := einolib.NewChatModel(ctx)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	_, _ = fmt.Fprint(os.Stdout, "[assistant] ")
	stream, err := cm.Stream(ctx, messages)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer stream.Close()

	for {
		frame, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			_, _ = fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		if frame != nil {
			_, _ = fmt.Fprint(os.Stdout, frame.Content)
		}
	}
	_, _ = fmt.Fprintln(os.Stdout)
}
