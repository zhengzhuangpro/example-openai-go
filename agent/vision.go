package agent

import (
	"context"
	"fmt"

	"github.com/openai/openai-go/v3"
)

// ChatVision 图片理解：多模态输入。
//
// 本示例学什么：
//   - 视觉模型的消息内容由多个 content part 拼成：
//     TextContentPart（文字提问）+ ImageContentPart（图片）
//   - 图片可以传公网 URL（本示例），也可以传 base64 data URL
//   - DeepSeek 使用deepseek-v4-flash-vision-exp https://api-docs.deepseek.com/zh-cn/news/news260821/
func (c *Client) ChatVision(ctx context.Context) error {
	// 想换图的话改这个 URL 即可（也可以换成 base64 data URL）
	imageURL := "http://e.hiphotos.baidu.com/image/pic/item/a1ec08fa513d2697e542494057fbb2fb4316d81e.jpg"

	resp, err := c.client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Model: openai.ChatModel("deepseek-v4-flash-vision-exp"),
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.UserMessage([]openai.ChatCompletionContentPartUnionParam{
				openai.TextContentPart("用中文描述这张图片，并猜猜它在哪里拍摄"),
				openai.ImageContentPart(openai.ChatCompletionContentPartImageImageURLParam{
					URL: imageURL,
				}),
			}),
		},
	})
	if err != nil {
		return fmt.Errorf("调用视觉模型失败: %w", err)
	}

	fmt.Println(resp.Choices[0].Message.Content)
	return nil
}
