package agent

import (
	"context"
	"fmt"

	"github.com/openai/openai-go/v3"
)

// ChatDeepSeekMultiTurn 多轮对话（DeepSeek 适配版）：手动拼接上下文。
//
// 本示例学什么：
// - DeepSeek 不支持 store/PreviousResponseID，多轮对话要自己维护 messages 切片
// - 每轮把「模型的回答」用 AssistantMessage 追加进历史，
//   模型才能「记得」之前说过什么
// - Chat Completions API 的结果在 resp.Choices[0].Message.Content
func (c *Client) ChatDeepSeekMultiTurn(ctx context.Context) error {
	messages := []openai.ChatCompletionMessageParamUnion{
		openai.UserMessage("法国的首都是哪里？"),
	}

	resp1, err := c.client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Model:    c.model,
		Messages: messages,
	})
	if err != nil {
		return err
	}

	answer1 := resp1.Choices[0].Message.Content
	fmt.Println("第一轮:", answer1)

	// 保存第一轮的上下文：模型回答 + 用户新问题
	messages = append(messages,
		openai.AssistantMessage(answer1),
		openai.UserMessage("那座城市有多少人口？"),
	)

	resp2, err := c.client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Model:    c.model,
		Messages: messages,
	})
	if err != nil {
		return err
	}

	fmt.Println("第二轮:", resp2.Choices[0].Message.Content)
	return nil
}
