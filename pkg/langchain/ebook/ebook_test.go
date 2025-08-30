/**
 * @author: dn-jinmin/dn-jinmin
 * @doc:
 */

package ebook

import (
	"context"
	"testing"

	"github.com/tmc/langchaingo/callbacks"
	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/llms/openai"
)

var (
	apiKey = "sk-eBf9hEicOVbobp2NbCRt67TsPB0fuhcl8iu09RpOKRJZj71s"
	url    = "https://api.openai-proxy.org/v1"
)

func getLLmOpenaiClient(t *testing.T, opts ...openai.Option) *openai.LLM {
	opts = append(opts, openai.WithBaseURL(url), openai.WithToken(apiKey))
	llm, err := openai.New(opts...)
	if err != nil {
		t.Fatal(err)
	}
	return llm
}

func TestNewEbook(t *testing.T) {
	llm := getLLmOpenaiClient(t, openai.WithModel("gpt-4o"), openai.WithCallback(callbacks.LogHandler{}))
	e := NewEbook(llm, "./upload/")
	res, err := chains.Call(context.Background(), e, map[string]any{
		"input": "我想写一篇关于上班族如何在上班的时候健身的文章",
	})
	t.Log(res, err)
}
