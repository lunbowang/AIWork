package callbackx

import (
	"context"

	"gitee.com/dn-jinmin/tlog"
	"github.com/tmc/langchaingo/callbacks"
	"github.com/tmc/langchaingo/llms"
)

// TitTokenHandle 是一个回调处理器，用于跟踪和记录LLM生成内容时使用的令牌(token)总数
type TitTokenHandle struct {
	callbacks.SimpleHandler
	logger tlog.Logger
}

// NewTitTokenHandle 创建一个新的TitTokenHandle实例
func NewTitTokenHandle(logger tlog.Logger) *TitTokenHandle {
	return &TitTokenHandle{
		SimpleHandler: callbacks.SimpleHandler{},
		logger:        logger,
	}
}

// HandleLLMGenerateContentEnd 处理LLM生成内容结束的事件，统计并记录总令牌数
func (l *TitTokenHandle) HandleLLMGenerateContentEnd(ctx context.Context, res *llms.ContentResponse) {
	var count int

	// 遍历所有生成的选择项(choices)，累加每个选择项的令牌数
	for i := range res.Choices {
		if v, ok := res.Choices[i].GenerationInfo["TotalTokens"]; ok {
			count += v.(int)
		}
	}

	if count == 0 {
		return
	}

	l.logger.InfofCtx(ctx, "TitTokenHandle:llm_generate_content_end", "count %d", count)
}
