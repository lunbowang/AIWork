package toolx

import (
	"ai/internal/svc"
	"ai/pkg/langchain"
	"context"
	"fmt"

	"github.com/tmc/langchaingo/callbacks"
	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/prompts"
)

const (
	// 修正模板：1. 动态接收 chatType 和 prompts；2. 确保 JSON 格式合法
	OUT_PROMPT_TEMPLATE = `<< instructions >>
1. First, output the thinking process and action (for agent parsing, cannot be omitted):
   Thought: I need to answer the user's question: {{.prompts}}
   Action: default  # Fixed value, do not change
   Action Input: {{.prompts}}  # Cannot be empty, use the user's question directly
2. Then output the final answer in valid JSON format (match the request structure):
   {"chatType": {{.chatType}}, "data": "{{.answer}}"}  # chatType is from request, data is your answer
3. "answer" must be a detailed response to {{.prompts}}, cannot be empty.
4. Do NOT output any extra content (no comments, blank lines, or markdown).`
)

type Empty struct {
	c        chains.Chain      // LLM 链
	callback callbacks.Handler // 回调（保留原逻辑）
}

// NewDefaultHandler：初始化默认处理器，模板改为动态接收 chatType 和 prompts
func NewDefaultHandler(svc *svc.ServiceContext) *Empty {
	// 主模板：整合用户输入、输出格式要求
	template := `You are an all-round assistant. The user's question is: 
<< User Input >>
{{.prompts}}

<< Output Requirements >>
{{.outPrompt}}`

	// 构建 PromptTemplate：输入变量包含 prompts/chatType/outPrompt（均动态传入）
	prompt := prompts.PromptTemplate{
		Template:       template,
		InputVariables: []string{"prompts", "chatType", "outPrompt"}, // 动态变量
		TemplateFormat: prompts.TemplateFormatGoTemplate,
		PartialVariables: map[string]any{
			"outPrompt": OUT_PROMPT_TEMPLATE, // 输出格式要求（固定，提前传入）
		},
	}

	return &Empty{
		c: chains.NewLLMChain(svc.LLMs, prompt, chains.WithCallback(svc.Callbacks)),
	}
}

// 保留原有的 Name/Description 方法（无需修改）
func (e *Empty) Name() string {
	return "default"
}

func (e *Empty) Description() string {
	return "这是一个默认的程序，在没有合适的选择时就选择它，在使用的时候请携带所有历史记录"
}

func (e *Empty) Call(ctx context.Context, input string) (string, error) {
	fmt.Println("empty -- call -- start --- ", input)
	out, err := chains.Predict(ctx, e.c, map[string]any{
		langchain.Input: input,
	})

	if err != nil {
		return "", err
	}

	return SuccessWithData + out + "\n\n\n", nil
}
