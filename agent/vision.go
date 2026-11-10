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
//   - DeepSeek 目前不支持视觉输入，这里切换到官方 client
//     （需要环境变量 OPENAI_API_KEY）
func (c *Client) ChatVision(ctx context.Context) error {
	client := NewOfficialClient()

	// 想换图的话改这个 URL 即可（也可以换成 base64 data URL）
	imageURL := "https://upload.wikimedia.org/wikipedia/commons/thumb/8/85/Tour_Eiffel_Wikimedia_Commons_%28cropped%29.jpg/320px-Tour_Eiffel_Wikimedia_Commons_%28cropped%29.jpg"

	resp, err := client.client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Model: client.model,
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
