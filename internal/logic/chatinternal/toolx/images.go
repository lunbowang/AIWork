package toolx

import (
	"ai/internal/domain"
	"ai/internal/svc"
	"ai/pkg/langchain/image"
	"context"
	"encoding/json"
	"fmt"
)

type Img struct {
	svc *svc.ServiceContext
	*image.Image
}

func NewImg(svc *svc.ServiceContext) *Img {
	return &Img{
		svc:   svc,
		Image: image.NewImage(svc.OpenaiClient, svc.Config.Upload.SavePath),
	}
}

func (i *Img) Name() string {
	return "img"
}

func (i *Img) Description() string {
	return "一个生成图片的接口，在用户需要生成图片的时候使用"
}

func (i *Img) Call(ctx context.Context, input string) (string, error) {
	fmt.Println("------- img call start ------- ")

	url, err := i.Gen(ctx, input, "")
	if err != nil {
		return "", err
	}

	data := domain.ChatResp{
		ChatType: domain.ImgAndText,
		Data: domain.ImgAndTextResp{
			Text: "",
			Url:  i.svc.Config.Upload.Host + "/" + url,
		},
	}

	body, err := json.Marshal(&data)
	if err != nil {
		return "", err
	}

	return SuccessWithData + string(body) + "\n\n\n", nil
}
