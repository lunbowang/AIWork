package voice

import (
	"context"
	"testing"

	openaiSdk "github.com/sashabaranov/go-openai"
)

func TestVoice_Transcriptions(t *testing.T) {
	var (
		url = "https://api.openai-proxy.org/v1"
		key = "sk-eBf9hEicOVbobp2NbCRt67TsPB0fuhcl8iu09RpOKRJZj71s"
	)
	openaiGPTCfg := openaiSdk.DefaultConfig(key)
	openaiGPTCfg.BaseURL = url
	openaiGPT := openaiSdk.NewClientWithConfig(openaiGPTCfg)

	t.Log(NewVoice(openaiGPT).Transcriptions(context.Background(), "./1.mp3", ""))
}
