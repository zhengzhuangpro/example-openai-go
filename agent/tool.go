package agent

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/example-openai-go/pkg/weather"
	"github.com/openai/openai-go/v3"
)

// ChatToolCalling Tool Calling 完整闭环：
// 模型决定调用工具 → 本地执行 → 结果回传 → 模型生成最终回答。
//
// 本示例学什么：
// - ChatCompletionFunctionTool 定义工具，JSON Schema 描述参数
// - 模型不直接给答案，而是在 ToolCalls 里给出「想调用的函数 + 参数」
// - 用 msg.ToParam() + openai.ToolMessage(结果, ID) 把工具结果回传
// - 最终回答同样可以走流式输出
func (c *Client) ChatToolCalling(ctx context.Context) error {
	// 1. 定义工具：告诉模型有哪些本地能力可用
	tools := []openai.ChatCompletionToolUnionParam{
		openai.ChatCompletionFunctionTool(openai.FunctionDefinitionParam{
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
		}),
	}

	history := []openai.ChatCompletionMessageParamUnion{
		openai.UserMessage("北京天气怎么样？"),
	}

	// 2. 第一轮请求：模型决定是否调用工具
	resp, err := c.client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Model:    c.model,
		Tools:    tools,
		Messages: history,
	})
	if err != nil {
		return err
	}

	msg := resp.Choices[0].Message

	// 注意：必须先判断非空再访问 ToolCalls。
	// 模型完全可能直接回答而不调工具，越界访问会 panic。
	if len(msg.ToolCalls) == 0 {
		fmt.Println("模型未调用工具，直接回答：", msg.Content)
		return nil
	}

	toolCall := msg.ToolCalls[0]
	fmt.Println("工具调用:", toolCall.Function.Name)

	// 3. 本地执行工具：解析模型给出的参数，调用真正的实现
	var args struct {
		City string `json:"city"`
	}
	if err := json.Unmarshal([]byte(toolCall.Function.Arguments), &args); err != nil {
		return fmt.Errorf("解析工具参数失败: %w", err)
	}

	result := weather.GetWeather(args.City, "celsius")

	// 4. 第二轮请求：「模型想调用工具」+「工具结果」一起回传
	stream := c.client.Chat.Completions.NewStreaming(ctx, openai.ChatCompletionNewParams{
		Model: c.model,
		Tools: tools,
		Messages: append(history,
			msg.ToParam(),
			openai.ToolMessage(result, toolCall.ID),
		),
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
