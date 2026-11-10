package agent

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/openai/openai-go/v3"
)

// ChatReasoner 推理模型思维链：deepseek-reasoner vs deepseek-chat。
//
// 本示例学什么：
//   - 同一个问题分别问普通模型和推理模型，观察回答差异
//   - reasoner 会先在 reasoning_content 里「自言自语」推理，
//     再给出最终答案——这是 DeepSeek 的非标字段，SDK 类型里没有，
//     要用 resp.RawJSON() 拿原始报文自己解析（见 extractReasoning）
//   - 模型名不必固定在 Client 上，可以按请求指定（这里就没用 c.model）
func (c *Client) ChatReasoner(ctx context.Context) error {
	question := "小明有 17 块糖，给了小红一半多一颗，又用剩下的糖换了 2 张卡片，" +
		"每张卡片花 4 块糖，最后他还剩几块糖？"

	// 1. 普通模型直接回答
	resp, err := c.client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Model:    openai.ChatModel("deepseek-chat"),
		Messages: []openai.ChatCompletionMessageParamUnion{openai.UserMessage(question)},
	})
	if err != nil {
		return err
	}
	fmt.Println("== deepseek-chat ==")
	fmt.Println(resp.Choices[0].Message.Content)
	fmt.Println()

	// 2. 推理模型：先出思维链，再出答案
	resp, err = c.client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Model:    openai.ChatModel("deepseek-reasoner"),
		Messages: []openai.ChatCompletionMessageParamUnion{openai.UserMessage(question)},
	})
	if err != nil {
		return err
	}

	msg := resp.Choices[0].Message

	fmt.Println("== deepseek-reasoner 思维链 ==")
	if reasoning, ok := extractReasoning(resp.RawJSON()); ok {
		fmt.Println(reasoning)
	} else {
		fmt.Println("（本模型未返回 reasoning_content）")
	}

	fmt.Println("\n== deepseek-reasoner 最终回答 ==")
	fmt.Println(msg.Content)

	return nil
}

// extractReasoning 从响应原始 JSON 里取 DeepSeek 的非标字段 reasoning_content。
// SDK 类型只定义了官方标准字段，第三方端点的扩展字段要用
// resp.RawJSON() 拿到原始报文后自己解析。
func extractReasoning(raw string) (string, bool) {
	var body struct {
		Choices []struct {
			Message struct {
				ReasoningContent string `json:"reasoning_content"`
			} `json:"message"`
		} `json:"choices"`
	}
	if err := json.Unmarshal([]byte(raw), &body); err != nil {
		return "", false
	}
	if len(body.Choices) == 0 || body.Choices[0].Message.ReasoningContent == "" {
		return "", false
	}
	return body.Choices[0].Message.ReasoningContent, true
}
