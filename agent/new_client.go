package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	getweater "github.com/example-openai-go/pkg/get_weater"
	"github.com/invopop/jsonschema"
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

// Tool Calling
func (c *NewClient) CreateChatCompletionToolRequest() (*responses.Response, error) {
	tools := []openai.ChatCompletionToolUnionParam{
		openai.ChatCompletionFunctionTool(openai.FunctionDefinitionParam{
			Name:        "get_weather",
			Description: openai.String("获取指定城市的天气信息"),
			Parameters: openai.FunctionParameters{
				"type": "object",
				"properties": map[string]any{
					"city": map[string]any{
						"type": "string",
					},
				},
				"required": []string{"city"},
			},
		}),
	}

	history := []openai.ChatCompletionMessageParamUnion{
		openai.UserMessage("北京天气怎么样？"),
	}

	// 第一次请求：模型决定调用哪个工具
	resp, err := c.client.Chat.Completions.New(c.ctx, openai.ChatCompletionNewParams{
		Model:    c.model,
		Tools:    tools,
		Messages: history,
	})
	if err != nil {
		return nil, err
	}

	msg := resp.Choices[0].Message

	fmt.Println("工具调用:", msg.ToolCalls[0].Function.Name)

	// 本地调用工具
	if len(msg.ToolCalls) > 0 {
		toolCall := msg.ToolCalls[0]

		var args struct {
			City string `json:"city"`
		}

		_ = json.Unmarshal([]byte(toolCall.Function.Arguments), &args)

		result := getweater.GetWeather(args.City, "celsius")

		// resp, err = c.client.Chat.Completions.New(c.ctx, openai.ChatCompletionNewParams{
		// 	Model: c.model,
		// 	Tools: tools,
		// 	Messages: append(history,
		// 		msg.ToParam(),
		// 		openai.ToolMessage(result, toolCall.ID),
		// 	),
		// })

		// 流式响应
		stream := c.client.Chat.Completions.NewStreaming(c.ctx, openai.ChatCompletionNewParams{
			Model: c.model,
			Tools: tools,
			Messages: append(history,
				msg.ToParam(),
				openai.ToolMessage(result, toolCall.ID),
			),
		})

		for stream.Next() {
			chunk := stream.Current()

			if len(chunk.Choices) > 0 {
				fmt.Print(chunk.Choices[0].Delta.Content)
			}
		}

		// 打印
		// fmt.Println(resp.Choices[0].Message.Content)

	}

	return nil, nil
}

type Origin struct {
	YearBuilt    int64  `json:"year_of_construction"`
	Organization string `json:"organization"`
}

type HistoricalComputer struct {
	Origin Origin `json:"origin" jsonSchema_description:"电脑的产地信息"`
	Model  string `json:"model" jsonSchema_description:"电脑的型号"`
	Year   int64  `json:"year" jsonSchema_description:"电脑的出厂年份"`
	Price  int64  `json:"price" jsonSchema_description:"电脑的出厂价格"`
}

func GenerateSchema[T any]() (map[string]any, error) {
	reflector := jsonschema.Reflector{
		AllowAdditionalProperties: false,
		DoNotReference:            true,
	}

	var v T

	data, err := json.Marshal(reflector.Reflect(v))
	if err != nil {
		return nil, err
	}

	var out map[string]any
	if err := json.Unmarshal(data, &out); err != nil {
		return nil, err
	}

	return out, nil
}

// JSON Schema 约束输出
func (c *NewClient) CreateChatCompletionJSONSchemaRequest() (*responses.Response, error) {
	schema, err := GenerateSchema[HistoricalComputer]()
	if err != nil {
		return nil, fmt.Errorf("生成 JSON Schema 失败: %w", err)
	}

	response, err := c.client.Responses.New(
		c.ctx,
		responses.ResponseNewParams{
			Model: c.model,
			Input: responses.ResponseNewParamsInputUnion{
				OfString: openai.String("什么计算机运行了第一个神经网络？"),
			},
			Text: responses.ResponseTextConfigParam{
				Format: responses.ResponseFormatTextConfigParamOfJSONSchema(
					"historical_computer",
					schema,
				),
			},
		},
	)
	if err != nil {
		return nil, fmt.Errorf("调用模型失败: %w", err)
	}

	var computer HistoricalComputer
	if err := json.Unmarshal(
		[]byte(response.OutputText()),
		&computer,
	); err != nil {
		return nil, fmt.Errorf("解析模型输出失败: %w", err)
	}

	fmt.Printf("JSON Schema 结果: %+v\n", computer)

	return response, nil
}
