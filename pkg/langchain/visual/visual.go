package visual

import (
	"ai/pkg/langchain"
	"context"
	"errors"
	"fmt"
	"github.com/tmc/langchaingo/callbacks"
	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/memory"
	"github.com/tmc/langchaingo/outputparser"
	"github.com/tmc/langchaingo/prompts"
	"github.com/tmc/langchaingo/schema"
)

type VisualTools interface {
	Namespace() string
	Destinations() string
	// inputs是针对于 具体的实现对象可以方便获取到数据
	Call(ctx context.Context, inputs map[string]string, options ...chains.ChainCallOption) (map[string]any, error)
}

// input -> db -> 数据的提取 -> 可视化

type VisualChains struct {
	Chain        chains.Chain
	callbacks    callbacks.Handler
	outputparser outputparser.Structured
	tools        map[string]VisualTools
	savePath     string
}

// option
func NewVisualChains(llm llms.Model, savePath string) *VisualChains {
	return &VisualChains{
		Chain: chains.NewLLMChain(llm, prompts.PromptTemplate{
			Template:       _defaultVisualParsePrompts,
			InputVariables: []string{langchain.Input},
			TemplateFormat: prompts.TemplateFormatGoTemplate,
			PartialVariables: map[string]any{
				_formatting: _outputParserVisual.GetFormatInstructions(),
			},
		}),
		callbacks:    nil,
		outputparser: _outputParserVisual,
		tools: map[string]VisualTools{
			"bar": NewBarTool(llm, savePath),
		},
		savePath: savePath,
	}
}

func (v *VisualChains) Call(ctx context.Context, inputs map[string]any, options ...chains.ChainCallOption) (map[string]any, error) {
	// 第一步获取分析的数据字段
	out, err := chains.Predict(ctx, v.Chain, inputs, options...)
	if err != nil {
		return inputs, err
	}

	// 解析信息
	info, err := v.extractInfo(out)
	if err != nil {
		return inputs, err
	}

	// 提取选择的方法
	next, ok := v.tools[info[_destinations]]
	if !ok {
		return inputs, err
	}

	return next.Call(ctx, info, options...)
}

func (v *VisualChains) extractInfo(out string) (map[string]string, error) {
	parseData, err := v.outputparser.Parse(out)
	if err != nil {
		return nil, err
	}

	info, ok := parseData.(map[string]string)
	if !ok {
		return nil, errors.New(fmt.Sprintf("VisualChains 不知道数据类型 %v", parseData))
	}
	return info, nil
}

func (v *VisualChains) GetMemory() schema.Memory {
	return memory.NewSimple()
}

func (v *VisualChains) GetInputKeys() []string {
	return []string{}
}

func (v *VisualChains) GetOutputKeys() []string {
	return []string{}
}
