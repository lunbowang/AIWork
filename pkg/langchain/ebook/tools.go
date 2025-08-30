package ebook

import (
	"ai/pkg/langchain"
	"context"
	"fmt"
	"strings"

	"github.com/tmc/langchaingo/callbacks"
	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/prompts"
)

// 用于生成文章主题、章节列表、摘要
type SubjectTool struct {
	chains.Chain
}

func NewSubjectTool(llm llms.Model) *SubjectTool {
	return &SubjectTool{
		Chain: chains.NewLLMChain(llm, prompts.NewPromptTemplate(SubjectToolPrompt, []string{langchain.Input}),
			chains.WithCallback(callbacks.LogHandler{})),
	}
}

func (s *SubjectTool) Call(ctx context.Context, input, filepath string) (string, error) {
	output, err := chains.Predict(ctx, s.Chain, map[string]any{
		langchain.Input: input,
	})
	if err != nil {
		return "", err
	}

	// 获取主题、摘要
	withoutMatchesStart := strings.Split(output, "章节：")
	if !(len(withoutMatchesStart) > 1) {
		return "", fmt.Errorf("SubjectTool 主题： 不存在信息 : %s", output)
	}
	withoutMatchesEnd := strings.Split(output, "摘要：")
	if !(len(withoutMatchesEnd) > 1) {
		return "", fmt.Errorf("SubjectTool 摘要： 不存在信息 : %s", output)
	}

	theme := strings.TrimSpace(withoutMatchesStart[0])
	abstract := strings.TrimSpace(withoutMatchesEnd[1])

	return output, SaveContent(fmt.Sprintf("%s\n摘要:\n%s\n", theme, abstract), filepath)
}

type ContentTool struct {
	chains.Chain
}

func NewContentTool(llm llms.Model) *ContentTool {
	return &ContentTool{
		Chain: chains.NewLLMChain(llm, prompts.NewPromptTemplate(ContentPrompt, []string{langchain.Input}),
			chains.WithCallback(callbacks.LogHandler{})),
	}
}

func (s *ContentTool) Call(ctx context.Context, input, filepath string) (string, error) {
	output, err := chains.Predict(ctx, s.Chain, map[string]any{
		langchain.Input: input,
	})
	if err != nil {
		return "", err
	}

	return output, SaveContent(output+"\n", filepath)
}
