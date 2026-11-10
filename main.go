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
	case "basic": // 基础调用：一次请求，OutputText 提取回复
		err = client.ChatBasic(ctx)
	case "multiturn": // 多轮对话：Store + PreviousResponseID 服务端续接
		err = client.ChatMultiTurn(ctx)
	case "deepseek": // 多轮对话：手动拼接 messages 上下文（DeepSeek 适配）
		err = client.ChatDeepSeekMultiTurn(ctx)
	case "conversations": // Conversations 会话管理 API（需官方端点）
		err = client.ChatConversations(ctx)
	case "stream": // 流式响应：增量打印生成内容
		err = client.ChatStream(ctx)
	case "tool": // Tool Calling：模型调用本地天气工具
		err = client.ChatToolCalling(ctx)
	case "tools": // 多工具路由：注册表 + 并行调用分发，Agent 雏形
		err = client.ChatMultiTool(ctx)
	case "reasoner": // 推理模型：deepseek-reasoner 思维链 vs 普通模型
		err = client.ChatReasoner(ctx)
	case "embeddings": // 语义搜索：向量化 + 余弦相似度（需官方端点）
		err = client.SearchEmbeddings(ctx)
	case "vision": // 图片理解：多模态 content parts（DeepSeek 视觉模型）
		err = client.ChatVision(ctx)
	case "schema": // JSON Schema：约束模型输出结构化数据
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
  tools          多工具路由：注册表 + 并行调用分发，Agent 雏形
  reasoner       推理模型：deepseek-reasoner 思维链 vs 普通模型
  embeddings     语义搜索：向量化 + 余弦相似度（需官方 OPENAI_API_KEY）
  vision         图片理解：多模态 content parts（deepseek-v4-flash-vision-exp）
  schema         JSON Schema：约束模型输出结构化数据

环境变量:
  DEEPSEEK_API_KEY  DeepSeek API Key（DeepSeek 示例必须）
  BAIDU_MAP_AK      百度地图 AK（tool/tools 子命令需要）
  OPENAI_API_KEY    OpenAI 官方 Key（embeddings/conversations 需要）`)
}
