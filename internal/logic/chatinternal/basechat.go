package chatinternal

import (
	"ai/internal/svc"
	"ai/pkg/langchain"
	"context"
	"strings"

	"github.com/tmc/langchaingo/agents"
	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/tools"
)

// baseChat 是基础聊天处理器结构体，封装了代理执行链
// 提供了通用的聊天处理功能，可被其他具体处理器继承
type baseChat struct {
	agentsChain chains.Chain
}

func NewBaseChat(svc *svc.ServiceContext, tools []tools.Tool) *baseChat {
	return &baseChat{
		// 创建代理执行链：
		// 1. 使用一次性代理(OneShotAgent)
		// 2. 传入LLM服务、可用工具和默认提示前缀
		agentsChain: agents.NewExecutor(agents.NewOneShotAgent(svc.LLMs, tools, agents.WithPromptPrefix(_defaultMrklPrefix))),
	}
}

// Chains 返回当前处理器的链式处理器
// 用于集成到更复杂的处理流程中
func (t *baseChat) Chains() chains.Chain {
	return chains.NewTransform(t.transform, nil, nil)
}

// transform 是核心转换方法，处理输入并返回输出
func (t *baseChat) transform(ctx context.Context, inputs map[string]any,
	opts ...chains.ChainCallOption) (map[string]any,
	error) {

	// 过滤输入：只保留字符串类型的输入值
	for s, a := range inputs {
		if _, ok := a.(string); !ok {
			delete(inputs, s)
		}
	}

	// 调用代理执行链处理输入
	outPut, err := t.agentsChain.Call(ctx, inputs, opts...)
	if err != nil {
		return nil, err
	}
	// 从输出中获取"output"字段的值
	v, ok := outPut["output"]
	if !ok {
		return outPut, nil
	}

	text := v.(string)

	// 处理JSON格式的输出：提取```json和```之间的内容
	withoutJSONStart := strings.Split(text, "```json")
	if !(len(withoutJSONStart) > 1) {
		// 如果没有JSON标记，直接返回原始输出
		return map[string]any{
			langchain.OutPut: v,
		}, err
	}

	withoutJSONEnd := strings.Split(withoutJSONStart[1], "```")
	if len(withoutJSONEnd) < 1 {
		// 如果JSON标记不完整，返回原始输出
		return map[string]any{
			langchain.OutPut: v,
		}, err
	}

	// 返回提取后的JSON内容作为输出
	return map[string]any{
		langchain.OutPut: withoutJSONEnd[0],
	}, nil
}
