package router

import (
	"fmt"
	"strings"

	"github.com/tmc/langchaingo/outputparser"
	"github.com/tmc/langchaingo/prompts"
)

const (
	_formatting   = "formatting"
	_destinations = "destinations"
	_input        = "input"
	_nextInput    = "next_inputs"

	MULTI_PROMPT_ROUTER_TEMPLATE = `Given a raw text input to a language model select the model prompt best suited for 
the input. You will be given the names of the available prompts and a description of 
what the prompt is best suited for. You may also revise the original input if you 
think that revising it will ultimately lead to a better response from the language 
model.

<< FORMATTING >>
Return a markdown code snippet with a JSON object formatted to look like:
{{.formatting}}

REMEMBER: "destination" MUST be one of the candidate prompt names specified below OR 
it can be "DEFAULT" if the input is not well suited for any of the candidate prompts.
REMEMBER: "next_inputs" can just be the original input if you don't think any 
modifications are needed.

<< CANDIDATE PROMPTS >>
{{.destinations}}

<< INPUT >>
{{.input}}
`
)

// 将 LLM（大语言模型）的输出解析为指定结构，确保输出包含目标问答系统名称和处理后的输入
var (
	// 全局结构化输出解析器实例
	_outputparser = outputparser.NewStructured([]outputparser.ResponseSchema{
		{
			// Name：解析字段名，对应 "_destinations"（目标问答系统标识）
			Name: _destinations,
			// Description：字段描述，明确该字段需填写"目标问答系统名称"或默认值"DEFAULT"
			Description: `name of the question answering system to use or "DEFAULT"`,
		}, {
			// Name：解析字段名，对应 "_nextInput"（处理后的输入）
			Name: _nextInput,
			// Description：字段描述，明确该字段需填写"原始输入的可能修改版本"
			Description: `a potentially modified version of the original input`,
		},
	})
)

// createPrompt：创建多问答系统路由的提示词模板
func createPrompt(handler []Handler) prompts.PromptTemplate {
	return prompts.PromptTemplate{
		// Template：提示词模板内容，使用预定义的多提示路由模板
		Template: MULTI_PROMPT_ROUTER_TEMPLATE,
		// InputVariables：模板所需的输入变量列表，此处仅需"原始输入（_input）"
		InputVariables: []string{_input},
		// TemplateFormat：模板格式，指定为 Go 模板语法（GoTemplate）
		TemplateFormat: prompts.TemplateFormatGoTemplate,
		// PartialVariables：部分预填充变量，提前注入模板中无需动态传入的固定内容
		PartialVariables: map[string]any{
			// _destinations：注入"目标系统列表"，通过 handlerDestinations 函数生成
			_destinations: handlerDestinations(handler),
			// _formatting：注入"输出格式指令"，通过输出解析器获取标准化格式要求
			_formatting: _outputparser.GetFormatInstructions(),
		},
	}
}

// handlerDestinations：将处理器列表转换为字符串格式的目标系统列表
func handlerDestinations(handler []Handler) string {
	var hs strings.Builder
	for _, h := range handler {
		hs.WriteString(fmt.Sprintf("- %s: %s\n", h.Name(), h.Description()))
	}

	return hs.String()
}
