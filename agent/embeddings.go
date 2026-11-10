package agent

import (
	"context"
	"fmt"
	"math"
	"sort"

	"github.com/openai/openai-go/v3"
)

// SearchEmbeddings 语义搜索：文本转向量 + 余弦相似度排序。
//
// 本示例学什么：
//   - Embeddings API 把文本变成高维向量，语义相近的文本向量距离近
//   - 一次请求可以传一组文本（InputUnion 传数组），批量向量化
//   - 余弦相似度衡量两个向量方向是否一致（与长度无关），
//     这是 RAG 检索的核心前置知识
//   - DeepSeek 不提供 Embeddings 端点，这里切换到官方 client
//     （需要环境变量 OPENAI_API_KEY）
func (c *Client) SearchEmbeddings(ctx context.Context) error {
	client := NewOfficialClient()

	// 迷你知识库：把几条 Go 知识向量化
	docs := []string{
		"goroutine 是 Go 并发的基本单元，由运行时调度，初始栈只有 2KB，可以轻松开出几十万个",
		"channel 用于 goroutine 之间传递数据，close 后接收端会读到零值，ok 返回 false",
		"slice 底层引用一段数组，append 超出容量时会分配新数组并拷贝，扩容策略从小到大翻倍递减",
		"interface 底层是 (动态类型, 动态值) 二元组，把具体值赋给接口时两者都非 nil",
		"GC 使用三色标记 + 混合写屏障，stop-the-world 只剩标记准备和收尾两个极短阶段",
	}
	query := "Go 里并发的数据竞争怎么避免？"

	fmt.Println("查询:", query)
	fmt.Println()

	// 知识库 + 查询一次性批量向量化（最后一条是查询本身）
	inputs := append(docs, query)
	resp, err := client.client.Embeddings.New(ctx, openai.EmbeddingNewParams{
		Model: openai.EmbeddingModelTextEmbedding3Small,
		Input: openai.EmbeddingNewParamsInputUnion{
			OfArrayOfStrings: inputs,
		},
	})
	if err != nil {
		return fmt.Errorf("调用 Embeddings 失败: %w", err)
	}

	// 响应按输入顺序返回向量，最后一条是查询向量
	vectors := make([][]float64, len(resp.Data))
	for _, d := range resp.Data {
		vectors[d.Index] = d.Embedding
	}
	queryVec := vectors[len(vectors)-1]

	// 计算每条知识与查询的余弦相似度，降序排序
	type hit struct {
		doc        string
		similarity float64
	}
	hits := make([]hit, 0, len(docs))
	for i, doc := range docs {
		hits = append(hits, hit{doc, cosineSimilarity(queryVec, vectors[i])})
	}
	sort.Slice(hits, func(i, j int) bool { return hits[i].similarity > hits[j].similarity })

	for i, h := range hits {
		fmt.Printf("第 %d 名 (相似度 %.4f): %s\n", i+1, h.similarity, h.doc)
	}

	return nil
}

// cosineSimilarity 余弦相似度：cos(θ) = A·B / (|A|×|B|)。
// 值域 [-1, 1]，越接近 1 表示两个向量方向越一致（语义越相近）。
func cosineSimilarity(a, b []float64) float64 {
	var dot, normA, normB float64
	for i := range a {
		dot += a[i] * b[i]
		normA += a[i] * a[i]
		normB += b[i] * b[i]
	}
	if normA == 0 || normB == 0 {
		return 0
	}
	return dot / (math.Sqrt(normA) * math.Sqrt(normB))
}
