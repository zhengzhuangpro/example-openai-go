// Package agent 封装 OpenAI Go SDK 的各种调用示例。
// 本 demo 通过 OpenAI 兼容端点接入 DeepSeek。
package agent

import (
	"net/http"
	"os"
	"time"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/shared"
)

// Client 封装 SDK 客户端与模型名。
type Client struct {
	client *openai.Client
	model  shared.ChatModel
}

// NewClient 创建客户端。
//
// 学什么：
// - option.WithAPIKey / WithBaseURL：OpenAI 兼容端点的标准接入方式，
//   DeepSeek、Moonshot、vLLM 等服务都这样接
// - option.WithHTTPClient：统一超时控制
// - option.WithMaxRetries：SDK 内置的网络重试
// - API Key 从环境变量 DEEPSEEK_API_KEY 读取，绝不硬编码进源码
//
// 注意：context 按惯例不存进结构体，而是随每次调用传入（见各示例方法）。
func NewClient() *Client {
	client := openai.NewClient(
		option.WithAPIKey(os.Getenv("DEEPSEEK_API_KEY")),
		option.WithBaseURL("https://api.deepseek.com"),
		option.WithHTTPClient(&http.Client{Timeout: 60 * time.Second}),
		option.WithMaxRetries(3),
	)
	return &Client{
		client: &client,
		model:  openai.ChatModel("deepseek-chat"), // 非官方模型名：用字符串转换构造
	}
}
