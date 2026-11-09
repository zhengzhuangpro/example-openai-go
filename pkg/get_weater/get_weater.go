package getweater

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// WeatherResult 工具返回给模型的结果结构
type WeatherResult struct {
	City        string  `json:"city"`
	Temperature float64 `json:"temperature"`
	Unit        string  `json:"unit"`
	Condition   string  `json:"condition"`
}

// getWeather 调用真实天气 API，返回 JSON 字符串
func GetWeather(city, unit string) string {
	// 可选参数兜底：模型不一定每次都传 unit
	if unit == "" {
		unit = "celsius"
	}

	// 1. 请求外部天气 API
	apiURL := fmt.Sprintf(
		"https://api.weather.example.com/v1/current?city=%s&unit=%s",
		url.QueryEscape(city), unit,
	)
	resp, err := http.Get(apiURL)
	if err != nil {
		return toolError("请求天气服务失败: " + err.Error())
	}
	defer resp.Body.Close()

	// 2. 解析 API 响应
	var apiResp struct {
		Temp      float64 `json:"temp"`
		Condition string  `json:"condition"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return toolError("解析天气数据失败: " + err.Error())
	}

	// 3. 序列化为 JSON 返回，模型解析更稳定
	result, _ := json.Marshal(WeatherResult{
		City: city, Temperature: apiResp.Temp, Unit: unit, Condition: apiResp.Condition,
	})
	return string(result)
}

// toolError 错误不走 panic，包装成 JSON 返回，让模型自行处理（重试或告知用户）
func toolError(msg string) string {
	b, _ := json.Marshal(map[string]string{"error": msg})
	return string(b)
}
