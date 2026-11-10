package agent

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/invopop/jsonschema"
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/responses"
)

// ChatJSONSchema 用 JSON Schema 约束模型输出结构化数据。
//
// 本示例学什么：
// - GenerateSchema[T]() 从 Go 结构体自动生成 JSON Schema（反射实现）
// - Text.Format 传入 schema 后，模型输出保证符合结构，可直接反序列化
// - 结构体 tag：json 控制字段名，jsonSchema_description 告诉模型字段含义
func (c *Client) ChatJSONSchema(ctx context.Context) error {
	schema, err := GenerateSchema[HistoricalComputer]()
	if err != nil {
		return fmt.Errorf("生成 JSON Schema 失败: %w", err)
	}

	response, err := c.client.Responses.New(ctx, responses.ResponseNewParams{
		Model: c.model,
		Input: responses.ResponseNewParamsInputUnion{
			OfString: openai.String("什么计算机运行了第一个神经网络？"),
		},
		Text: responses.ResponseTextConfigParam{
			Format: responses.ResponseFormatTextConfigParamOfJSONSchema(
				"historical_computer",
				schema,
			),
		},
	})
	if err != nil {
		return fmt.Errorf("调用模型失败: %w", err)
	}

	// 输出已被 schema 约束，可以放心直接反序列化成结构体
	var computer HistoricalComputer
	if err := json.Unmarshal([]byte(response.OutputText()), &computer); err != nil {
		return fmt.Errorf("解析模型输出失败: %w", err)
	}

	fmt.Printf("JSON Schema 结果: %+v\n", computer)
	return nil
}

type Origin struct {
	YearBuilt    int64  `json:"year_of_construction"`
	Organization string `json:"organization"`
}

type HistoricalComputer struct {
	Origin Origin `json:"origin" jsonSchema_description:"电脑的产地信息"`
	Model  string `json:"model" jsonSchema_description:"电脑的型号"`
	Year   int64  `json:"year" jsonSchema_description:"电脑的出厂年份"`
	Price  int64  `json:"price" jsonSchema_description:"电脑的出厂价格"`
}

// GenerateSchema 从任意结构体类型生成 JSON Schema（map 形式，可直接传给 SDK）。
//
// 学什么：
// - jsonschema.Reflector 靠反射读结构体 tag 生成 schema
// - DoNotReference: true 把嵌套结构内联展开（SDK 的 Format 参数要求平铺）
// - AllowAdditionalProperties: false 禁止模型输出多余字段
func GenerateSchema[T any]() (map[string]any, error) {
	reflector := jsonschema.Reflector{
		AllowAdditionalProperties: false,
		DoNotReference:            true,
	}

	var v T

	data, err := json.Marshal(reflector.Reflect(v))
	if err != nil {
		return nil, err
	}

	var out map[string]any
	if err := json.Unmarshal(data, &out); err != nil {
		return nil, err
	}

	return out, nil
}
