package analysis

import (
	"ai/internal/domain"
	"ai/internal/svc"
	"ai/pkg/langchain"
	"ai/pkg/langchain/visual"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/sashabaranov/go-openai"

	"github.com/tmc/langchaingo/callbacks"
	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/outputparser"
	"github.com/tmc/langchaingo/prompts"
)

type File struct {
	svc          *svc.ServiceContext
	c            chains.Chain
	outputParser outputparser.Structured
	visual       chains.Chain
	callbacks    callbacks.Handler
}

func NewFile(svc *svc.ServiceContext) *File {
	return &File{
		svc:          svc,
		outputParser: _defaultFileOutputParser,
		c: chains.NewLLMChain(svc.LLMs, prompts.PromptTemplate{
			Template:       _defaultFileQuery,
			InputVariables: []string{langchain.Input},
			TemplateFormat: prompts.TemplateFormatGoTemplate,
			PartialVariables: map[string]any{
				"formatting": _defaultFileOutputParser.GetFormatInstructions(),
			},
		}),
		visual: visual.NewVisualChains(svc.LLMs, svc.Config.Upload.SavePath),
	}
}

func (d *File) Name() string {
	return "file"
}

func (d *File) Description() string {
	return `根据历史对话和用户的信息表明指定对某一个文件分析的时候使用`
}

func (d *File) Chains() chains.Chain {
	return chains.NewTransform(d.transform, nil, nil)
}

func (d *File) transform(ctx context.Context, inputs map[string]any, options ...chains.ChainCallOption) (map[string]any,
	error) {
	// ？
	fmt.Println("analysis --- file --- start")

	// 获取文件的名称、地址、文件id
	out, err := chains.Predict(ctx, d.c, inputs, options...)
	if err != nil {
		return nil, err
	}

	// 对参数解析
	query, err := d.extractQuery(ctx, out)
	if err != nil {
		return nil, err
	}
	fmt.Printf("analysis file d.extractQuery query %v\n", query)

	fileId := query[FileId]
	isExistFileId := true // 避免重复设置文件id到记忆机制里
	if len(fileId) == 0 || fileId == "1" {
		// 就要上传文件,并返回文件id
		fileId, err = d.uploadFile(ctx, query)
		if err != nil {
			return nil, err
		}
		fmt.Printf("analysis file upload file query %v , fileid %v \n", query, fileId)
		isExistFileId = false
	}

	// 数据分析
	res, err := d.svc.AliProxyOpenai.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model: "qwen-long",
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: query[System] + "如果是进行数据分析，请不要输出具体细节",
			}, {
				Role:    openai.ChatMessageRoleSystem,
				Content: "fileid://" + fileId,
			}, {
				Role:    openai.ChatMessageRoleUser,
				Content: inputs[langchain.Input].(string),
			},
		},
	})
	if err != nil {
		return nil, err
	}
	fmt.Printf("analysis file CreateChatCompletion res %v\n", res)

	// 进行数据的可视化
	v, err := chains.Call(ctx, d.visual, map[string]any{
		langchain.Input: res.Choices[0].Message.Content,
	})

	var resp domain.ChatResp
	if err != nil {
		fmt.Println("visual : ", err)
		resp = domain.ChatResp{
			ChatType: 0,
			Data:     res.Choices[0].Message.Content,
		}
	} else {
		resp = domain.ChatResp{
			ChatType: domain.ImgAndText,
			Data: domain.ImgAndTextResp{
				Text: res.Choices[0].Message.Content,
				Url:  d.svc.Config.Upload.Host + v[visual.ImgUrl].(string),
			},
		}
	}
	body, err := json.Marshal(&resp)
	if err != nil {
		return nil, err
	}

	if isExistFileId {
		return map[string]any{
			langchain.OutPut: string(body),
		}, nil
	}

	return map[string]any{
		langchain.OutPut: string(body),
		langchain.Input:  fmt.Sprintf("%s fileId : %s", query[FileName], fileId),
	}, nil
}

func (f *File) extractQuery(ctx context.Context, text string) (map[string]string, error) {
	res, err := f.outputParser.Parse(text)
	if err != nil {
		return nil, err
	}

	// callback

	data, ok := res.(map[string]string)
	if !ok {
		return nil, errors.New("不知道什么类型")
	}

	return data, nil
}

func (f *File) uploadFile(ctx context.Context, fileInfo map[string]string) (string, error) {
	if _, err := os.Stat(fileInfo[FilePath]); err != nil {
		return "", errors.New("不存在文件")
	}

	file, err := f.svc.AliProxyOpenai.CreateFile(ctx, openai.FileRequest{
		FileName: fileInfo[FileName],
		FilePath: fileInfo[FilePath],
		Purpose:  "file-extract",
	})
	if err != nil {
		return "", err
	}

	fmt.Println("upload find id : ", file.ID)

	// 对文件的删除
	go func(fid string) {
		time.Sleep(time.Hour)
		f.svc.AliProxyOpenai.DeleteFile(context.Background(), fid)
	}(file.ID)

	return file.ID, err
}
