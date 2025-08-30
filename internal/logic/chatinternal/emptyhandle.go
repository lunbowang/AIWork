package chatinternal

import (
	"ai/internal/domain"
	"ai/internal/svc"
	"ai/pkg/langchain"
	"fmt"

	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/prompts"
)

type DefaultHandler struct {
	svc *svc.ServiceContext
	c   chains.Chain
}

// NewDefaultHandler 创建默认处理器实例
func NewDefaultHandler(svc *svc.ServiceContext) *DefaultHandler {
	template := "you are an all-round assistant,please help me answer this question: \n\n<< input >>\n{{.input}}"

	prompt := prompts.PromptTemplate{
		Template:       BASE_PROMPAT_TEMPLATE + template + "\n\n" + OUT_PROMPT_TEMPLATE, // 拼接角色模板与输出格式模板
		InputVariables: []string{langchain.Input},                                       // 模板输入变量：仅需"input"（用户问题）
		TemplateFormat: prompts.TemplateFormatGoTemplate,                                // 模板格式：Go模板语法
		// 部分预填充变量：无需动态传入，提前注入模板的固定内容
		PartialVariables: map[string]any{
			"chatType": fmt.Sprintf("%d", domain.DefaultHandler),
			"data":     "solution",
		},
	}

	// 初始化LLM链：关联LLM模型、提示词模板，并绑定服务上下文的回调处理器
	return &DefaultHandler{
		svc: svc,
		c: chains.NewLLMChain(
			svc.LLMs,                           // 服务上下文提供的LLM实例
			prompt,                             // 构建好的提示词模板
			chains.WithCallback(svc.Callbacks), // 绑定全局回调（如日志、令牌统计）
		),
	}
}

func (d *DefaultHandler) Name() string {
	return "default"
}

func (d *DefaultHandler) Description() string {
	return "suitable for answering multiple questions"
}

// Chains 返回默认处理器关联的LLM链，供路由系统调用以处理用户问题
func (d *DefaultHandler) Chains() chains.Chain {
	return d.c
}
