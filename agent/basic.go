package agent

import (
	"context"
	"fmt"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/responses"
)

// ChatBasic 基础调用：一次请求，拿到完整回复。
//
// 本示例学什么：
// - Responses API 最小调用：只需 Model + Input
// - Input 用 ResponseNewParamsInputUnion 包装，OfString 传纯文本
// - resp.OutputText() 直接提取最终文本，不用自己遍历 output 结构
func (c *Client) ChatBasic(ctx context.Context) error {
	resp, err := c.client.Responses.New(ctx, responses.ResponseNewParams{
		Model: c.model,
		Input: responses.ResponseNewParamsInputUnion{
			OfString: openai.String("用一句话解释 goroutine"),
		},
	})
	if err != nil {
		return err
	}

	fmt.Println(resp.OutputText())
	return nil
}
