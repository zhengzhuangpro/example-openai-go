package agent

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

// ChatErrors 错误处理：区分 API 错误和网络错误。
//
// 本示例学什么：
//   - SDK 的 API 错误是 *openai.Error，用 errors.As 断言后能拿到
//     StatusCode / Code / Message —— 限流(429)、鉴权(401)、参数错(400)
//     的应对方式各不相同
//   - 网络错误（连不上、超时）不是 API 错误，断言失败走普通 error 分支
//   - 生产建议：429 指数退避重试（client.go 的 WithMaxRetries 已内置）、
//     4xx 是请求本身有问题（改代码，别重试）、5xx 可以重试
func (c *Client) ChatErrors(ctx context.Context) error {
	// 错误 1：模型名不存在 → API 返回 4xx
	_, err := c.client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Model:    openai.ChatModel("not-exist-model"),
		Messages: []openai.ChatCompletionMessageParamUnion{openai.UserMessage("你好")},
	})
	describeError(err)

	// 错误 2：连不上的端点 → 网络错误，不是 API 错误
	badClient := openai.NewClient(
		option.WithBaseURL("https://invalid.example-openai-go.dev"),
		option.WithHTTPClient(&http.Client{Timeout: 3 * time.Second}),
	)
	_, err = badClient.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Model:    c.model,
		Messages: []openai.ChatCompletionMessageParamUnion{openai.UserMessage("你好")},
	})
	describeError(err)

	return nil
}

// describeError 把 SDK 返回的错误分类打印。
func describeError(err error) {
	fmt.Println("---- 错误 ----")

	var apiErr *openai.Error
	if errors.As(err, &apiErr) {
		fmt.Println("类型:   API 错误")
		fmt.Println("状态码:", apiErr.StatusCode)
		fmt.Println("错误码:", apiErr.Code)
		fmt.Println("信息:  ", apiErr.Message)
		return
	}

	fmt.Println("类型: 其他错误（网络/本地），errors.As 断言失败")
	fmt.Println("信息:  ", err)
}
