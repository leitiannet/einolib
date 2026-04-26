// 提供模型与智能体的构造器，用于组合 einolib 与实现包中的常见默认配置
package builder

import (
	"context"
	"fmt"

	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/components/model"
	"github.com/leitiannet/einolib"
	libchatmodel "github.com/leitiannet/einolib/agents/chatmodel"
	libdeep "github.com/leitiannet/einolib/agents/deep"
	libfilesystem "github.com/leitiannet/einolib/middlewares/filesystem"
	_ "github.com/leitiannet/einolib/models"
	_ "github.com/leitiannet/einolib/operators"
)

// 本地模型的默认参数
const (
	DefaultLocalModelType    = einolib.ModelTypeOllama
	DefaultLocalModelID      = "qwen2:7b"
	DefaultLocalModelBaseURL = "http://localhost:11434"
)

// 聊天模型构造器
type ModelBuilder struct {
	modelOptions []einolib.ModelOption
}

func NewModelBuilder(modelOptions ...einolib.ModelOption) *ModelBuilder {
	return &ModelBuilder{modelOptions: modelOptions}
}

func (b *ModelBuilder) Build(ctx context.Context) (model.ToolCallingChatModel, error) {
	modelOptions := []einolib.ModelOption{
		einolib.WithModelType(DefaultLocalModelType),
		einolib.WithModelDescription("model from builder"),
		einolib.WithModelID(DefaultLocalModelID),
		einolib.WithModelBaseURL(DefaultLocalModelBaseURL),
	}
	modelOptions = append(modelOptions, b.modelOptions...)
	return einolib.NewChatModel(ctx, modelOptions...)
}

// 智能体构造器
type AgentBuilder struct {
	agentType        einolib.AgentType
	agentName        string
	agentDescription string
	instruction      string
	modelOptions     []einolib.ModelOption
}

func NewAgentBuilder(agentType einolib.AgentType, agentName, agentDescription string) *AgentBuilder {
	return &AgentBuilder{
		agentType:        agentType,
		agentName:        agentName,
		agentDescription: agentDescription,
	}
}

func (b *AgentBuilder) WithInstruction(instruction string) *AgentBuilder {
	b.instruction = instruction
	return b
}

func (b *AgentBuilder) WithModelOptions(modelOptions ...einolib.ModelOption) *AgentBuilder {
	b.modelOptions = append(b.modelOptions, modelOptions...)
	return b
}

func (b *AgentBuilder) Build(ctx context.Context) (adk.Agent, error) {
	if agentType := b.agentType; agentType != einolib.AgentTypeChatModel && agentType != einolib.AgentTypeDeep {
		return nil, fmt.Errorf("invalid agent type: %s", agentType)
	}
	chatModel, err := CreateChatModel(ctx, b.modelOptions)
	if err != nil {
		return nil, err
	}
	localOperator, err := CreateLocalOperator(ctx)
	if err != nil {
		return nil, err
	}
	agentOptions := []einolib.AgentOption{
		einolib.WithAgentType(b.agentType),
		einolib.WithAgentName(b.agentName),
		einolib.WithAgentDescription(b.agentDescription),
	}
	if b.agentType == einolib.AgentTypeChatModel {
		fileSystemMiddleware, err := CreateFileSystemMiddleware(ctx, localOperator)
		if err != nil {
			return nil, err
		}
		chatModelAgentConfig := libchatmodel.NewChatModelAgentConfig(
			libchatmodel.WithModel(chatModel),
			libchatmodel.WithHandlers(fileSystemMiddleware),
			libchatmodel.WithInstruction(b.instruction),
		)
		agentOptions = append(agentOptions, einolib.WithAgentComponentConfig(chatModelAgentConfig))
	} else {
		deepAgentConfig := libdeep.NewDeepAgentConfig(
			libdeep.WithChatModel(chatModel),
			libdeep.WithBackend(localOperator),
			libdeep.WithShell(localOperator),
			libdeep.WithStreamingShell(localOperator),
			libdeep.WithInstruction(b.instruction),
		)
		agentOptions = append(agentOptions, einolib.WithAgentComponentConfig(deepAgentConfig))
	}
	return einolib.NewAgent(ctx, agentOptions...)
}

func CreateChatModel(ctx context.Context, modelOptions []einolib.ModelOption) (model.ToolCallingChatModel, error) {
	return NewModelBuilder(modelOptions...).Build(ctx)
}

func CreateLocalOperator(ctx context.Context) (einolib.Operator, error) {
	operatorOptions := []einolib.OperatorOption{
		einolib.WithOperatorType(einolib.OperatorTypeLocal),
		einolib.WithOperatorDescription("filesystem operator from builder"),
	}
	return einolib.NewOperator(ctx, operatorOptions...)
}

func CreateFileSystemMiddleware(ctx context.Context, operator einolib.Operator) (adk.ChatModelAgentMiddleware, error) {
	middlewareOptions := []einolib.MiddlewareOption{
		einolib.WithMiddlewareType(einolib.MiddlewareTypeFileSystem),
		einolib.WithMiddlewareDescription("filesystem middleware from builder"),
		einolib.WithMiddlewareComponentConfig(libfilesystem.NewFileSystemMiddlewareConfig(
			libfilesystem.WithOperator(operator),
		)),
	}
	return einolib.NewMiddleware(ctx, middlewareOptions...)
}
