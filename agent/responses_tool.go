package agent

import (
	"context"
	"fmt"
	"time"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/responses"
)

// ChatResponsesTool Responses API 版 Tool Calling。
//
// 本示例学什么：
//   - 和 tool.go（Chat Completions）同样的功能，换 Responses API 实现，
//     对照感受两种 API 的差异
//   - 工具定义：responses.FunctionToolParam，挂在 ResponseNewParams.Tools
//   - 响应输出：resp.Output 是条目列表，type == "function_call" 的条目
//     带着模型想调的工具名和参数（CallID / Name / Arguments）
//   - 回传：Input 换成「条目列表」，把用户消息、模型的调用意图、
//     工具输出按顺序放回去（OfFunctionCall / OfFunctionCallOutput）
func (c *Client) ChatResponsesTool(ctx context.Context) error {
	tools := []responses.ToolUnionParam{{
		OfFunction: &responses.FunctionToolParam{
			Name:        "get_time",
			Description: openai.String("获取当前的本地时间"),
			Parameters:  map[string]any{"type": "object", "properties": map[string]any{}},
		},
	}}

	// 第一轮：模型决定调用工具
	resp, err := c.client.Responses.New(ctx, responses.ResponseNewParams{
		Model: c.model,
		Input: responses.ResponseNewParamsInputUnion{OfString: openai.String("现在几点了？")},
		Tools: tools,
	})
	if err != nil {
		return err
	}

	// 从输出条目里找 function_call
	var callID, callName, callArgs string
	for _, item := range resp.Output {
		if item.Type == "function_call" {
			callID, callName, callArgs = item.CallID, item.Name, item.Arguments.OfString
		}
	}
	if callID == "" {
		fmt.Println("模型未调用工具：", resp.OutputText())
		return nil
	}

	fmt.Println("工具调用:", callName, "参数:", callArgs)

	// 本地执行工具
	result := time.Now().Format(`{"time":"2006-01-02 15:04:05"}`)

	// 第二轮：回传完整条目链，拿最终回答
	resp, err = c.client.Responses.New(ctx, responses.ResponseNewParams{
		Model: c.model,
		Tools: tools,
		Input: responses.ResponseNewParamsInputUnion{
			OfInputItemList: responses.ResponseInputParam{
				{OfMessage: &responses.EasyInputMessageParam{
					Role: responses.EasyInputMessageRoleUser,
					Content: responses.EasyInputMessageContentUnionParam{
						OfString: openai.String("现在几点了？"),
					},
				}},
				{OfFunctionCall: &responses.ResponseFunctionToolCallParam{
					CallID: callID, Name: callName, Arguments: callArgs,
				}},
				{OfFunctionCallOutput: &responses.ResponseInputItemFunctionCallOutputParam{
					CallID: openai.String(callID),
					Output: responses.ResponseInputItemFunctionCallOutputOutputUnionParam{
						OfString: openai.String(result),
					},
				}},
			},
		},
	})
	if err != nil {
		return err
	}

	fmt.Println("最终回答：", resp.OutputText())
	return nil
}
