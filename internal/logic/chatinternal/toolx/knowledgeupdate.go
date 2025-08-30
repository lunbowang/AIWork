package toolx

import (
	"ai/internal/svc"
	"ai/pkg/langchain/outputparserx"
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/tmc/langchaingo/embeddings"

	"github.com/tmc/langchaingo/textsplitter"

	"github.com/tmc/langchaingo/callbacks"
	"github.com/tmc/langchaingo/documentloaders"
	"github.com/tmc/langchaingo/vectorstores/redisvector"
)

type KnowledgeUpdate struct {
	svc          *svc.ServiceContext
	Callback     callbacks.Handler
	outPutParser outputparserx.Structured
	store        *redisvector.Store
}

func NewKnowledgeUpdate(svc *svc.ServiceContext) *KnowledgeUpdate {
	return &KnowledgeUpdate{
		svc:      svc,
		Callback: svc.Callbacks,
		outPutParser: outputparserx.NewStructured([]outputparserx.ResponseSchema{
			{
				Name:        "path",
				Description: "the path to file",
			}, {
				Name:        "name",
				Description: "the name to file",
			}, {
				Name:        "time",
				Description: "file update time",
			},
		}),
	}
}

func (k KnowledgeUpdate) Name() string {
	return "knowledge_update"
}

func (k KnowledgeUpdate) Description() string {
	return `
a knowledge base update interface.
use when you need to update knowledge base content.
your output should be in the following json format.
` + k.outPutParser.GetFormatInstructions()
}

// Call 执行知识库更新操作
func (k *KnowledgeUpdate) Call(ctx context.Context, input string) (string, error) {
	// 权限验证：检查当前上下文是否有权限执行知识库更新操作
	if err := k.svc.Auth(ctx); err != nil {
		return "", err
	}

	var data any // 用于存储解析后的输入数据
	// 解析输入字符串，将其转换为结构化数据
	data, err := k.outPutParser.Parse(input)
	if err != nil {
		// ```json str
		// 解析失败时，尝试直接将输入作为JSON字符串解析
		t := make(map[string]any)
		if err := json.Unmarshal([]byte(input), &t); err != nil {
			return "", err // JSON解析也失败，返回错误
		}

		return "", err // 原始解析失败，返回错误
	}

	// 将解析后的数据断言为map类型，便于获取文件路径
	file := data.(map[string]any)

	// 打开指定路径的文件（从解析后的数据中获取路径）
	f, err := os.Open(fmt.Sprintf("%v", file["path"]))
	if err != nil {
		return "", err
	}
	// 获取文件信息（如大小、修改时间等）
	finfo, err := f.Stat()
	if err != nil {
		return "", err
	}

	// 创建PDF文档加载器，用于读取PDF内容
	p := documentloaders.NewPDF(f, finfo.Size())
	// 加载PDF并分割为文本块：
	// - 每个块大小为200个字符
	// - 块之间重叠1个字符（确保上下文连贯性）
	chunkedDocuments, err := p.LoadAndSplit(ctx, textsplitter.NewRecursiveCharacter(textsplitter.WithChunkSize(200),
		textsplitter.WithChunkOverlap(1)))
	if err != nil {
		return "", err // 文档加载或分割失败，返回错误
	}

	// 如果知识库存储实例未初始化，则创建它
	if k.store == nil {
		k.store, err = getKnowledgeStore(ctx, k.svc)
		if err != nil {
			return "", err
		}
	}

	// 将分割后的文本块添加到知识库中
	_, err = k.store.AddDocuments(ctx, chunkedDocuments)
	if err != nil {
		return "", err
	}

	return Success, nil
}

// getKnowledgeStore 创建并返回Redis向量存储实例（用于存储知识库向量数据）
func getKnowledgeStore(ctx context.Context, svc *svc.ServiceContext) (*redisvector.Store, error) {
	// 创建嵌入器（用于将文本转换为向量表示）
	embedder, err := embeddings.NewEmbedder(svc.LLMs)
	if err != nil {
		return nil, err
	}

	// 创建并返回Redis向量存储实例：
	// - 使用前面创建的嵌入器
	// - 连接到配置中指定的Redis服务
	// - 指定索引名称为"knowledge"，并设置为可创建（不存在则创建）
	return redisvector.New(ctx, redisvector.WithEmbedder(embedder), redisvector.WithConnectionURL("redis://"+svc.Config.
		Redis.
		Addr), redisvector.WithIndexName("knowledge", true))
}
