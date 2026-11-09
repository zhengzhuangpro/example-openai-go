package agent

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/conversations"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/packages/pagination"
	"github.com/openai/openai-go/v3/responses"
	"github.com/openai/openai-go/v3/shared"
)

type NewClient struct {
	ctx    context.Context
	client *openai.Client
	model  shared.ChatModel
}

func NewNewClient(ctx context.Context) *NewClient {
	// 创建一个新的OpenAI客户端
	client := openai.NewClient(
		option.WithAPIKey(os.Getenv("DEEPSEEK_API_KEY")),
		option.WithBaseURL("https://api.deepseek.com"),
		option.WithHTTPClient(&http.Client{Timeout: 60 * time.Second}),
		option.WithMaxRetries(3), // 设置最大重试次数
	)
	return &NewClient{
		ctx:    ctx,
		client: &client,
		model:  openai.ChatModel("deepseek-chat"), // 非官方模型：字符串转换构造器
	}
}

// 创建一个新的聊天完成请求, 基本调用
func (c *NewClient) CreateChatCompletionRequest(messages responses.ResponseNewParamsInputUnion) (*responses.Response, error) {
	resp, err := c.client.Responses.New(c.ctx, responses.ResponseNewParams{
		Model: c.model,
		Input: messages,
	})
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// 多轮对话
func (c *NewClient) CreateChatCompletionRequestWithMessages() (*responses.Response, error) {
	resp1, err := c.client.Responses.New(c.ctx, responses.ResponseNewParams{
		Model: c.model,
		Store: openai.Bool(true), // 存储 Response，支持后续引用, 多轮对话
		Input: responses.ResponseNewParamsInputUnion{OfString: openai.String("法国的首都是哪里？")},
	})
	if err != nil {
		return nil, err
	}

	fmt.Println("第一轮:", resp1.OutputText())

	resp2, err := c.client.Responses.New(c.ctx, responses.ResponseNewParams{
		Model:              c.model,
		Store:              openai.Bool(true),       // 开启多轮对话
		PreviousResponseID: openai.String(resp1.ID), // 接着上一轮
		Input:              responses.ResponseNewParamsInputUnion{OfString: openai.String("那座城市有多少人口？")},
	})
	if err != nil {
		return nil, err
	}
	fmt.Println("第二轮:", resp2.OutputText())
	fmt.Println("最新 ResponseID:", resp2.ID)
	return resp2, nil
}

// 多轮会话deepseek不支持store，响应中恒为 store: false，这里写一个deepseek适配的例子
func (c *NewClient) CreateChatCompletionRequestWithMessagesDeepSeek() (*string, error) {
	messages := []openai.ChatCompletionMessageParamUnion{
		openai.UserMessage("法国的首都是哪里？"),
	}

	resp1, err := c.client.Chat.Completions.New(c.ctx, openai.ChatCompletionNewParams{
		Model:    c.model,
		Messages: messages,
	})
	if err != nil {
		return nil, err
	}

	fmt.Println("第一轮:", resp1.Choices[0].Message.Content)
	answer1 := resp1.Choices[0].Message.Content

	// 保存第一轮的上下文
	messages = append(messages,
		openai.AssistantMessage(answer1),
		openai.UserMessage("那座城市有多少人口？"),
	)

	resp2, err := c.client.Chat.Completions.New(c.ctx, openai.ChatCompletionNewParams{
		Model:    c.model,
		Messages: messages,
	})
	if err != nil {
		return nil, err
	}
	fmt.Println("第二轮:", resp2.Choices[0].Message.Content)
	return &resp2.Choices[0].Message.Content, nil
}

// Conversations API（会话管理） deepseek不支持
func (c *NewClient) CreateChatCompletionRequestWithConversations() (*responses.Response, error) {
	var (
		conv  *conversations.Conversation
		resp  *responses.Response
		err   error
		items *pagination.ConversationCursorPage[conversations.ConversationItemUnion]
	)

	if conv, err = c.client.Conversations.New(c.ctx, conversations.ConversationNewParams{}); err != nil {
		return nil, err
	}

	if resp, err = c.client.Responses.New(c.ctx, responses.ResponseNewParams{
		Model: c.model,
		Input: responses.ResponseNewParamsInputUnion{OfString: openai.String("记住我喜欢的颜色是蓝色")},
		Conversation: responses.ResponseNewParamsConversationUnion{
			OfConversationObject: &responses.ResponseConversationParam{ID: conv.ID},
		},
	}); err != nil {
		return nil, err
	}

	if items, err = c.client.Conversations.Items.List(c.ctx, conv.ID, conversations.ItemListParams{}); err != nil {
		return nil, err
	}

	fmt.Println("会话条目数:", len(items.Data))
	return resp, nil
}

// 流式响应
func (c *NewClient) CreateChatCompletionStreamRequest() (*responses.Response, error) {
	stream := c.client.Responses.NewStreaming(c.ctx, responses.ResponseNewParams{
		Model: c.model,
		Input: responses.ResponseNewParamsInputUnion{
			OfString: openai.String("写一首关于编程的俳句"),
		},
	})

	for stream.Next() {
		fmt.Print(stream.Current().Delta) // 文本增量 内容
	}

	return nil, nil
}
