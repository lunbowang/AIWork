package image

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/sashabaranov/go-openai"
)

type Image struct {
	openaiClient *openai.Client
	savePath     string
}

func NewImage(client *openai.Client, savePath string) *Image {
	return &Image{
		openaiClient: client,
		savePath:     savePath,
	}
}

func (i *Image) Gen(ctx context.Context, prompt, size string) (string, error) {
	if len(size) == 0 {
		size = "1024x1024"
	}

	res, err := i.openaiClient.CreateImage(ctx, openai.ImageRequest{
		Prompt: prompt,
		Model:  "dall-e-3",
		N:      1,
		Size:   size,
	})
	if err != nil {
		return "", err
	}

	// 保存图
	filePath := fmt.Sprintf("%s_img_%d.png", i.savePath, time.Now().Unix())
	if err := downloadImage(res.Data[0].URL, filePath); err != nil {
		return "", err
	}
	return filePath, nil
}
func downloadImage(url string, filepath string) error {
	// 发送 HTTP GET 请求
	response, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to download image: %v", err)
	}
	defer response.Body.Close()

	// 检查 HTTP 响应状态码
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download image: %s", response.Status)
	}

	// 创建文件
	out, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}
	defer out.Close()

	// 将图片内容写入文件
	_, err = io.Copy(out, response.Body)
	if err != nil {
		return fmt.Errorf("failed to save image: %v", err)
	}

	return nil
}
