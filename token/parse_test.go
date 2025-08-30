package token

import (
	"fmt"
	"net/http"
	"testing"
	"time"
)

// TestSecretKey 测试用的密钥，用于JWT令牌的生成和验
var TestSecretKey = "LunBoWang"

func TestGenToken(t *testing.T) {
	now := time.Now().Unix()
	t.Log(GetJwtToken(TestSecretKey, now, 60*60*60, "1"))
}

func TestVerifyJWTToken(t *testing.T) {

	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTc2MTg0NjMsImlhdCI6MTcxNzQwMjQ2MywiaW1vb2MuY29tIjoiMSJ9.anVWrthElU1ZS34UcFpE380aSvp30KtWq1_CIl6YnKo"

	r := &http.Request{
		Header: make(http.Header),
	}
	r.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token))

	t.Log(VerifyJWTToken(TestSecretKey, r))
}
