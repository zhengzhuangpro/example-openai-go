package agent

import (
	"context"
	"fmt"

	"github.com/openai/openai-go/v3"
)

// ChatParams 参数与用量：temperature 对生成的影响 + token 统计。
//
// 本示例学什么：
//   - temperature 越低输出越确定（适合代码、抽取），
//     越高越发散（适合创作）；0 基本每次都一样
//   - resp.Usage 是 token 用量（输入/输出/总计），
//     计费核算和上下文窗口管理都靠它
func (c *Client) ChatParams(ctx context.Context) error {
	prompt := "用一个词形容秋天"

	for _, temp := range []float64{0, 2} {
		resp, err := c.client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
			Model:       c.model,
			Messages:    []openai.ChatCompletionMessageParamUnion{openai.UserMessage(prompt)},
			Temperature: openai.Float(temp),
		})
		if err != nil {
			return err
		}

		fmt.Printf("temperature=%.1f → %s\n", temp, resp.Choices[0].Message.Content)
		fmt.Printf("  用量: 输入 %d + 输出 %d = 总计 %d tokens\n",
			resp.Usage.PromptTokens, resp.Usage.CompletionTokens, resp.Usage.TotalTokens)
	}

	return nil
}
