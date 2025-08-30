package chatinternal

import (
	"ai/token"
	"context"
	"testing"

	"github.com/tmc/langchaingo/chains"
)

func TestChatLogHandle_transform(t *testing.T) {
	chat := NewChatLogHandle(svcTest)
	tokenStr := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJMdW5Cb1dhbmciOiI2OGE4NjI4OWY5Yjg4ZjQ4ODI0ODQxOTQiLCJleHAiOjE3NjQ4NTQ3NDUsImlhdCI6MTc1NjIxNDc0NX0.a8WvQWit-SUYF-YWkJ9rIH15jaseOqOhle9V41oMtXg"
	ctx := context.Background()
	ctx = context.WithValue(ctx, token.Authorization, tokenStr)

	res, err := chat.transform(ctx, map[string]any{
		"relationId": "5JRwHwo9IVp5ffyovmlQrX",
		"startTime":  1756130224,
		"endTime":    1756130266,
	}, chains.WithCallback(svcTest.Callbacks))
	if err != nil {
		t.Fatal(err)
	}

	t.Log(res)
}
