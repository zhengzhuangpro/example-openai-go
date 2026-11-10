package agent

import (
	"context"
	"fmt"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/responses"
)

// ChatMultiTurn 多轮对话：让服务端帮你记上下文。
//
// 本示例学什么：
// - Store: true 把 Response 存在服务端，拿到可引用的 ID
// - PreviousResponseID 指定「接着哪一轮说」，第二轮只需传新问题，
//   历史上下文由服务端拼接
// - 注意：DeepSeek 不支持 store（响应中恒为 store: false），
//   手动拼接上下文的方案见 deepseek.go
func (c *Client) ChatMultiTurn(ctx context.Context) error {
	resp1, err := c.client.Responses.New(ctx, responses.ResponseNewParams{
		Model: c.model,
		Store: openai.Bool(true),
		Input: responses.ResponseNewParamsInputUnion{OfString: openai.String("法国的首都是哪里？")},
	})
	if err != nil {
		return err
	}
	fmt.Println("第一轮:", resp1.OutputText())

	resp2, err := c.client.Responses.New(ctx, responses.ResponseNewParams{
		Model:              c.model,
		Store:              openai.Bool(true),       // 开启多轮对话
		PreviousResponseID: openai.String(resp1.ID), // 接着上一轮
		Input:              responses.ResponseNewParamsInputUnion{OfString: openai.String("那座城市有多少人口？")},
	})
	if err != nil {
		return err
	}
	fmt.Println("第二轮:", resp2.OutputText())
	fmt.Println("最新 ResponseID:", resp2.ID)
	return nil
}
