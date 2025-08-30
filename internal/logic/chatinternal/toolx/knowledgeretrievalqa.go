package toolx

import (
	"ai/internal/svc"
	"context"

	"github.com/tmc/langchaingo/vectorstores"

	"github.com/tmc/langchaingo/callbacks"
	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/vectorstores/redisvector"
)

type KnowledgeRetrievalQA struct {
	svc      *svc.ServiceContext
	Callback callbacks.Handler
	store    *redisvector.Store
	qa       chains.Chain
}

func NewKnowledgeRetrievalQA(svc *svc.ServiceContext) *KnowledgeRetrievalQA {
	return &KnowledgeRetrievalQA{svc: svc}
}

func (k *KnowledgeRetrievalQA) Name() string {
	return "knowledge_retrieval_qa"
}

func (k *KnowledgeRetrievalQA) Description() string {
	return `a knowledge retrieval interface.
use it when you need to inquire about work-related policies, such as employee manuals, attendance rules, etc.
keep Chinese output.`
}

// Call 执行知识库检索问答操作
// ctx: 上下文对象，用于传递超时控制和请求上下文
// input: 用户的查询字符串（问题）
// 返回值: 检索到的答案字符串和可能的错误
func (k *KnowledgeRetrievalQA) Call(ctx context.Context, input string) (string, error) {
	var err error
	// 初始化问答链（如果尚未初始化）
	if k.qa == nil {
		// 获取知识库存储实例（连接到Redis向量存储）
		k.store, err = getKnowledgeStore(ctx, k.svc)
		if err != nil {
			return "", err
		}

		// 创建检索式问答链：
		// - 使用服务中配置的LLM（大语言模型）作为回答生成器
		// - 将向量存储转换为检索器，设置每次检索返回1个最相关的结果
		k.qa = chains.NewRetrievalQAFromLLM(k.svc.LLMs, vectorstores.ToRetriever(k.store, 1))
	}

	// 调用问答链处理用户查询
	// 传入查询参数（"query"字段对应用户输入）
	res, err := chains.Predict(ctx, k.qa, map[string]any{
		"query": input,
	})
	if err != nil {
		return "", err
	}

	return `The following are the consultation results. When outputting, please output the results directly, do not make summaries, keep them in Chinese, and only do the original output:
\n` + res, nil
}
