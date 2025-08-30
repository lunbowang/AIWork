package chatinternal

import (
	"ai/pkg/langchain"
	"ai/token"
	"context"
	"testing"

	"github.com/tmc/langchaingo/chains"
)

func TestApproval(t *testing.T) {
	tokenStr := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJMdW5Cb1dhbmciOiI2OGE4NjI4OWY5Yjg4ZjQ4ODI0ODQxOTQiLCJleHAiOjE3NjQ4NTQ3NDUsImlhdCI6MTc1NjIxNDc0NX0.a8WvQWit-SUYF-YWkJ9rIH15jaseOqOhle9V41oMtXg"

	ctx := context.Background()
	ctx = context.WithValue(ctx, token.Authorization, tokenStr)

	chat := NewApprovalHandle(svcTest)
	res, err := chat.baseChat.transform(ctx, map[string]any{
		langchain.Input: "提交一个明天上午请假审批",
		"history":       "",
	}, chains.WithCallback(svcTest.Callbacks))

	if err != nil {
		t.Fatal(err)
	}

	t.Log(res)
}
