package agent

import (
	"context"
	"fmt"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/conversations"
	"github.com/openai/openai-go/v3/responses"
)

// ChatConversations Conversations API：服务端会话管理。
//
// 本示例学什么：
// - Conversations.New 先创建一个「会话容器」
// - Response 里通过 Conversation 字段挂到会话上，服务端自动归档上下文
// - Conversations.Items.List 可以列出会话里的全部条目
// - 注意：DeepSeek 不支持该 API，此示例需切到官方端点运行
func (c *Client) ChatConversations(ctx context.Context) error {
	conv, err := c.client.Conversations.New(ctx, conversations.ConversationNewParams{})
	if err != nil {
		return err
	}

	resp, err := c.client.Responses.New(ctx, responses.ResponseNewParams{
		Model: c.model,
		Input: responses.ResponseNewParamsInputUnion{OfString: openai.String("记住我喜欢的颜色是蓝色")},
		Conversation: responses.ResponseNewParamsConversationUnion{
			OfConversationObject: &responses.ResponseConversationParam{ID: conv.ID},
		},
	})
	if err != nil {
		return err
	}

	items, err := c.client.Conversations.Items.List(ctx, conv.ID, conversations.ItemListParams{})
	if err != nil {
		return err
	}

	fmt.Println("模型回复:", resp.OutputText())
	fmt.Println("会话条目数:", len(items.Data))
	return nil
}
