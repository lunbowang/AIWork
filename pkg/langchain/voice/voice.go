package voice

import (
	"context"

	"github.com/sashabaranov/go-openai"
)

type Voice struct {
	openai *openai.Client
}

func NewVoice(client *openai.Client) *Voice {
	return &Voice{
		openai: client,
	}
}

func (v *Voice) Transcriptions(ctx context.Context, filepath, prompts string) (string, error) {
	res, err := v.openai.CreateTranscription(ctx, openai.AudioRequest{
		Model:    "whisper-1",
		FilePath: filepath,
		Prompt:   prompts,
	})

	if err != nil {
		return "", err
	}

	return res.Text, nil
}
