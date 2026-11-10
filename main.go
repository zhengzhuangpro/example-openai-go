// example-openai-go：OpenAI Go SDK 学习示例集（通过兼容端点接入 DeepSeek）。
//
// 用法：go run . <子命令>，不带参数查看全部示例。
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/example-openai-go/agent"
)

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(1)
	}

	client := agent.NewClient()
	ctx := context.Background()

	var err error
	switch os.Args[1] {
	case "basic":
		err = client.ChatBasic(ctx)
	case "multiturn":
		err = client.ChatMultiTurn(ctx)
	case "deepseek":
		err = client.ChatDeepSeekMultiTurn(ctx)
	case "conversations":
		err = client.ChatConversations(ctx)
	case "stream":
		err = client.ChatStream(ctx)
	case "tool":
		err = client.ChatToolCalling(ctx)
	case "schema":
		err = client.ChatJSONSchema(ctx)
	case "help", "-h", "--help":
		usage()
		return
	default:
		fmt.Fprintln(os.Stderr, "未知子命令:", os.Args[1])
		usage()
		os.Exit(1)
	}

	if err != nil {
		log.Fatal(err)
	}
}

func usage() {
	fmt.Println(`OpenAI Go SDK 学习示例（DeepSeek 兼容端点）

用法: go run . <子命令>

子命令:
  basic          基础调用：一次请求，OutputText 提取回复
  multiturn      多轮对话：Store + PreviousResponseID 服务端续接
  deepseek       多轮对话：手动拼接 messages 上下文（DeepSeek 适配）
  conversations  Conversations 会话管理 API（DeepSeek 不支持，需官方端点）
  stream         流式响应：增量打印生成内容
  tool           Tool Calling：模型调用本地天气工具
  schema         JSON Schema：约束模型输出结构化数据

环境变量:
  DEEPSEEK_API_KEY  DeepSeek API Key（必须）
  BAIDU_MAP_AK      百度地图 AK（tool 子命令需要）`)
}
