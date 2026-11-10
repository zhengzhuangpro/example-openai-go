// Package weather 封装百度地图天气 API，作为 Tool Calling 的本地工具实现。
package weather

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
)

// GetWeather 调用真实天气 API，返回 JSON 字符串给模型解析。
//
// 学什么：
//   - 工具的返回值是给「模型」看的，返回 JSON 字符串最稳定
//   - 出错不走 panic：包装成 {"error": "..."} 返回，让模型自行处理
//     （重试或告诉用户），这也避免了程序崩溃
//   - AK 这类凭证一律从环境变量读取，绝不硬编码进源码
func GetWeather(city, unit string) string {
	// 可选参数兜底：模型不一定每次都传 unit
	if unit == "" {
		unit = "celsius"
	}

	fmt.Println("调用外部天气 API，获取城市：", city, "单位：", unit)

	ak := os.Getenv("BAIDU_MAP_AK")
	if ak == "" {
		return toolError("缺少环境变量 BAIDU_MAP_AK，请在百度地图开放平台申请后再试")
	}

	// 百度天气 API 按行政区划查询（district 传全名，如「北京市」），
	// 这里的 city 参数直接透传给 district
	params := url.Values{
		"district":  []string{city},
		"data_type": []string{"all"},
		"ak":        []string{ak},
	}
	apiURL := "https://api.map.baidu.com/weather/v1/?" + params.Encode()

	resp, err := http.Get(apiURL)
	if err != nil {
		return toolError("请求天气服务失败: " + err.Error())
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body) // io.ReadAll，ioutil 已废弃
	if err != nil {
		return toolError("读取天气数据失败: " + err.Error())
	}

	return string(body)
}

// toolError 错误不走 panic，包装成 JSON 返回，让模型自行处理（重试或告知用户）
func toolError(msg string) string {
	b, _ := json.Marshal(map[string]string{"error": msg})
	return string(b)
}
