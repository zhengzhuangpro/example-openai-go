package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/example-openai-go/pkg/weather"
	"github.com/openai/openai-go/v3"
)

// localTool 一个本地工具 = 参数 schema + 真正的执行函数。
// 注册进 map，模型喊什么名字就路由到什么实现。
type localTool struct {
	def  openai.FunctionDefinitionParam
	exec func(args json.RawMessage) string
}

// ChatMultiTool 多工具路由：真实 Agent 的雏形。
//
// 本示例学什么：
//   - 把「工具定义」和「工具实现」绑在一起注册成表，
//     新增工具只需往 map 里加一项，分发逻辑完全不用改
//   - 模型可能一次并行返回多个 ToolCalls（问时间 + 问天气），
//     要逐个执行、逐个用 ToolMessage 回传，一个都不能漏
//   - 生产级 Agent 还要把「请求-调用工具-回传」包成 while 循环，
//     直到模型不再要求调用工具为止（本示例演示一轮，保持简洁）
func (c *Client) ChatMultiTool(ctx context.Context) error {
	registry := map[string]localTool{
		"get_weather": {
			def: openai.FunctionDefinitionParam{
				Name:        "get_weather",
				Description: openai.String("获取指定城市的天气信息"),
				Parameters: openai.FunctionParameters{
					"type": "object",
					"properties": map[string]any{
						"city": map[string]any{
							"type":        "string",
							"description": "城市名，需匹配行政区划全名，例如：北京市",
						},
					},
					"required": []string{"city"},
				},
			},
			exec: func(args json.RawMessage) string {
				var a struct {
					City string `json:"city"`
				}
				if err := json.Unmarshal(args, &a); err != nil {
					return `{"error":"解析 city 参数失败"}`
				}
				return weather.GetWeather(a.City, "celsius")
			},
		},
		"get_time": {
			def: openai.FunctionDefinitionParam{
				Name:        "get_time",
				Description: openai.String("获取当前的本地时间"),
				Parameters: openai.FunctionParameters{
					"type":       "object",
					"properties": map[string]any{}, // 无参数工具
				},
			},
			exec: func(args json.RawMessage) string {
				return time.Now().Format(`{"time":"2006-01-02 15:04:05"}`) // 直接拼进格式串
			},
		},
	}

	// 从注册表收集工具定义，传给模型
	tools := make([]openai.ChatCompletionToolUnionParam, 0, len(registry))
	for _, t := range registry {
		tools = append(tools, openai.ChatCompletionFunctionTool(t.def))
	}

	history := []openai.ChatCompletionMessageParamUnion{
		openai.UserMessage("现在几点了？北京今天天气怎么样？"),
	}

	// 第一轮：模型决定调用哪些工具
	resp, err := c.client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Model:    c.model,
		Tools:    tools,
		Messages: history,
	})
	if err != nil {
		return err
	}

	msg := resp.Choices[0].Message
	if len(msg.ToolCalls) == 0 {
		fmt.Println("模型未调用工具，直接回答：", msg.Content)
		return nil
	}

	// 回传时先带上「模型的调用意图」，再逐个补上工具结果
	messages := append(history, msg.ToParam())

	// 模型可能并行发起多个工具调用，逐个路由执行
	for _, toolCall := range msg.ToolCalls {
		fmt.Println("工具调用:", toolCall.Function.Name, "参数:", toolCall.Function.Arguments)

		tool, ok := registry[toolCall.Function.Name]
		if !ok {
			// 模型幻觉出不存在的工具：明确回传错误，让它自己纠正
			messages = append(messages, openai.ToolMessage(
				`{"error":"不存在的工具"}`, toolCall.ID,
			))
			continue
		}

		result := tool.exec(json.RawMessage(toolCall.Function.Arguments))
		messages = append(messages, openai.ToolMessage(result, toolCall.ID))
	}

	// 第二轮：模型拿到全部工具结果，流式生成最终回答
	stream := c.client.Chat.Completions.NewStreaming(ctx, openai.ChatCompletionNewParams{
		Model:    c.model,
		Tools:    tools,
		Messages: messages,
	})

	fmt.Println("最终回答：")
	for stream.Next() {
		if chunk := stream.Current(); len(chunk.Choices) > 0 {
			fmt.Print(chunk.Choices[0].Delta.Content)
		}
	}
	if err := stream.Err(); err != nil {
		return fmt.Errorf("流式响应失败: %w", err)
	}
	fmt.Println()

	return nil
}
