package main

import (
	"context"
	"fmt"

	"github.com/example-openai-go/agent"
)

func main() {
	fmt.Println("Hello, World!")

	ctx := context.Background()

	// 创建一个新的OpenAI客户端
	client := agent.NewNewClient(ctx)

	// 创建一个新的聊天请求
	// resp, err := client.CreateChatCompletionRequest(responses.ResponseNewParamsInputUnion{
	// 	OfString: openai.String("用一句话解释 goroutine"),
	// })
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println(resp.OutputText())

	// 多轮对话
	// resp, err := client.CreateChatCompletionRequestWithMessages()
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println(resp.OutputText())

	// 多轮对话，deepseek支持
	// resp, err := client.CreateChatCompletionRequestWithMessagesDeepSeek()
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println(resp)

	// 会话管理
	// resp, err := client.CreateChatCompletionRequestWithConversations()
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println(resp)

	// 流式响应
	// resp, err := client.CreateChatCompletionStreamRequest()
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println(resp)

	// Tool Calling
	// resp, err := client.CreateChatCompletionToolRequest()
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println(resp)

	// JSON Schema
	_, err := client.CreateChatCompletionJSONSchemaRequest()
	if err != nil {
		panic(err)
	}
	// fmt.Println(resp)
}
