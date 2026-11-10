# example-openai-go

OpenAI Go SDK 学习示例集，通过 OpenAI 兼容端点接入 [DeepSeek](https://api.deepseek.com)。每个示例对应 `agent/` 下的一个文件，跑一个命令学一个知识点。

## 准备

```bash
export DEEPSEEK_API_KEY=sk-xxx   # 必须
export BAIDU_MAP_AK=xxx          # 仅 tool 子命令需要（百度地图开放平台申请）
```

## 示例清单

| 子命令 | 文件 | 学什么 |
| --- | --- | --- |
| `basic` | `agent/basic.go` | Responses API 最小调用，`OutputText()` 提取回复 |
| `multiturn` | `agent/multiturn.go` | `Store` + `PreviousResponseID` 让服务端记上下文 |
| `deepseek` | `agent/deepseek.go` | DeepSeek 不支持 store，手动维护 messages 切片 |
| `conversations` | `agent/conversations.go` | Conversations 会话管理 API（DeepSeek 不支持） |
| `stream` | `agent/stream.go` | 流式响应，`stream.Next()` 增量消费 + `stream.Err()` |
| `tool` | `agent/tool.go` | Tool Calling 闭环：模型决定调用 → 本地执行 → 结果回传 |
| `schema` | `agent/schema.go` | JSON Schema 约束输出，结构体反射生成 schema |

## 运行

```bash
go run . basic
go run . tool
go run . schema
```

不带参数运行可查看完整用法帮助。

## 目录结构

```
main.go              # 子命令入口
agent/               # 每个示例一个文件
  client.go          #   Client 封装与构造
pkg/weather/         # Tool Calling 的本地工具实现（百度天气 API）
```
