package main

import (
	"context"
	"fmt"

	"github.com/sashabaranov/go-openai"
)

var (
	apiKey = "sk-21fc35649d0e4063a8000a2a5f8cf2c0"
	url    = "https://dashscope.aliyuncs.com/compatible-mode/v1"
)

func main() {
	// 配置
	gptCfg := openai.DefaultConfig(apiKey)
	gptCfg.BaseURL = url
	ctx := context.Background()
	// 客户端
	client := openai.NewClientWithConfig(gptCfg)

	// 上传文件
	file, err := client.CreateFile(ctx, openai.FileRequest{
		FileName: "财务报表",
		FilePath: "1.xlsx",
		Purpose:  "file-extract",
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("file id : ", file.ID)

	//id := "file-fe-9iqYGBxa0mUkhPpZdCwkFsoo"

	// 根据文件提问
	res, err := client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model: "qwen-long",
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: "你是一个财务分析助理",
			}, {
				Role:    openai.ChatMessageRoleSystem,
				Content: "fileid://" + file.ID,
			}, {
				Role:    openai.ChatMessageRoleUser,
				Content: "请帮我分析出待发生、已完成、进行中分别的总金额是多少 禁止具体细节输出",
			},
		},
	})
	fmt.Println("res: ", res, " err ", err)
}
