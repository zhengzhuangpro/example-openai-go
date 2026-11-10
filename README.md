# example-openai-go

OpenAI Go SDK 学习示例集，通过 OpenAI 兼容端点接入 [DeepSeek](https://api.deepseek.com)。每个示例对应 `agent/` 下的一个文件，跑一个命令学一个知识点。

## 准备

```bash
export DEEPSEEK_API_KEY=sk-xxx   # DeepSeek 示例必须
export BAIDU_MAP_AK=xxx          # tool/tools 子命令需要（百度地图开放平台申请）
export OPENAI_API_KEY=sk-xxx     # embeddings/vision/conversations 需要（官方端点）
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
| `tools` | `agent/tools.go` | 多工具注册表 + map 路由分发，模型并行调用多个工具 |
| `reasoner` | `agent/reasoner.go` | deepseek-reasoner 思维链（ExtraFields 取非标字段）vs 普通模型 |
| `embeddings` | `agent/embeddings.go` | 文本向量化 + 余弦相似度语义搜索（需官方 key） |
| `vision` | `agent/vision.go` | TextContentPart + ImageContentPart 多模态输入（deepseek-v4-flash-vision-exp） |
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
