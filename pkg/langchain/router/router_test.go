package router

import (
	"ai/pkg/langchain/callbackx"
	"context"
	"testing"

	"gitee.com/dn-jinmin/tlog"

	"github.com/tmc/langchaingo/callbacks"
	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"
	"github.com/tmc/langchaingo/prompts"
)

var (
	apiKey = "sk-eBf9hEicOVbobp2NbCRt67TsPB0fuhcl8iu09RpOKRJZj71s"
	url    = "https://api.openai-proxy.org/v1"
)

// getLLmOpenaiClient 创建并返回OpenAI LLM客户端实例
func getLLmOpenaiClient(t *testing.T, opts ...openai.Option) *openai.LLM {
	// 合并默认配置（API地址和密钥）与自定义选项
	opts = append(opts, openai.WithBaseURL(url), openai.WithToken(apiKey))
	// 创建OpenAI客户端
	llm, err := openai.New(opts...)
	if err != nil {
		t.Fatal(err) // 初始化失败时终止测试
	}
	return llm
}

// TestRouter 测试路由功能，验证是否能正确将问题路由到对应的处理器
func TestRouter(t *testing.T) {
	logger := tlog.NewLogger(tlog.WithLogWriteLimit(1))

	// 初始化回调处理器组合（当前为空，可添加日志、令牌统计等回调）
	callback := callbacks.CombiningHandler{
		Callbacks: []callbacks.Handler{
			callbackx.NewLogHandler(logger),
			callbackx.NewTitTokenHandle(logger),
		},
	}

	// 创建带有回调的OpenAI LLM客户端
	llms := getLLmOpenaiClient(t, openai.WithCallback(callback))

	// 定义处理器列表：包含游泳和篮球相关的处理器
	handlers := []Handler{
		NewSwimmingHandler(llms),
		NewBasketballHandler(llms),
	}

	// 创建路由实例，传入LLM、处理器列表和回调配置
	router := NewRouter(llms, handlers, Withcallback(callback))

	res, err := chains.Call(tlog.TraceStart(context.Background()), router, map[string]any{
		"input": "请问游泳的类型有哪些",
	}, chains.WithCallback(callback))
	//res, err := chains.Call(context.Background(), router, map[string]any{
	//	"input": "乒乓球怎么打",
	//}, chains.WithCallback(callback))

	t.Log(res)
	t.Log(err)

}

type SwimmingHandler struct {
	c chains.Chain
}

func NewSwimmingHandler(llms llms.Model) *SwimmingHandler {
	return &SwimmingHandler{c: chains.NewLLMChain(
		llms,
		prompts.NewPromptTemplate(
			"你是一个资深游泳教练，精通游泳相关的所有知识, 请回答下面的问题：\n{{.input}}", []string{"input"},
		),
	)}
}

func (h SwimmingHandler) Name() string {
	return "swimming"
}

func (h SwimmingHandler) Description() string {
	return "适合回答游泳相关的知识"
}
func (h SwimmingHandler) Chains() chains.Chain {
	return h.c
}

type BasketballHandler struct {
	c chains.Chain
}

func NewBasketballHandler(llms llms.Model) *BasketballHandler {
	return &BasketballHandler{c: chains.NewLLMChain(
		llms,
		prompts.NewPromptTemplate(
			"你是一个资深篮球教练，精通篮球相关的所有知识, 请回答下面的问题：\n{{.input}}", []string{"input"},
		),
	)}
}

func (h BasketballHandler) Name() string {
	return "basketball"
}

func (h BasketballHandler) Description() string {
	return "适合回答篮球相关的知识"
}
func (h BasketballHandler) Chains() chains.Chain {
	return h.c
}
