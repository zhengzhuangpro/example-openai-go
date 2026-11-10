package agent

import (
	"context"
	"fmt"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/responses"
)

// ChatStream 流式响应：内容边生成边输出。
//
// 本示例学什么：
// - NewStreaming 返回一个流对象，for stream.Next() 逐块消费
// - stream.Current() 拿到当前增量块，Delta 是这一段新生成的文本
// - stream.Next() 返回 false 不代表一定正常结束，
//   必须再检查 stream.Err() —— 否则网络错误会被静默吞掉
func (c *Client) ChatStream(ctx context.Context) error {
	stream := c.client.Responses.NewStreaming(ctx, responses.ResponseNewParams{
		Model: c.model,
		Input: responses.ResponseNewParamsInputUnion{
			OfString: openai.String("写一首关于编程的俳句"),
		},
	})

	for stream.Next() {
		fmt.Print(stream.Current().Delta) // 文本增量
	}
	if err := stream.Err(); err != nil {
		return fmt.Errorf("流式响应失败: %w", err)
	}
	fmt.Println()
	return nil
}
