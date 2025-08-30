package chatinternal

import (
	"ai/pkg/langchain"
	"context"
	"testing"
)

func TestNewDefaultHandler(t *testing.T) {
	defaultHandler := NewDefaultHandler(svcTest)
	ctx := context.Background()

	res, err := defaultHandler.baseChat.transform(ctx, map[string]any{
		langchain.Input: "请帮我生成一张星空的图片",
		"history":       "",
	})
	if err != nil {
		t.Fatal(err)
	}
	//res, err := defaultHandler.baseChat.transform(ctx, map[string]any{
	//	langchain.Input: "你可以做什么",
	//	"history":       "AI,你好我是AI小慧，可以生成图片和文章",
	//})
	//if err != nil {
	//	t.Fatal(err)
	//}
	t.Log(res)

}
