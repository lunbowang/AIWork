package toolx

import (
	"ai/internal/domain"
	"ai/pkg/httpx"
	"encoding/json"
	"errors"
)

var (
	Success = `executes successfully. `

	SuccessWithData = `
executes successfully.
After you have determined the final answer, you should not make changes to the content, do not summarize the content, do not output your thoughts, and only keep the output of the original results.
Keep the output in json format as follows.\n
`
)

// ResParser 解析AI聊天接口的响应结果，并根据不同的聊天类型进行处理
func ResParser(v []byte, chatType domain.AiChatType, err error) (string, error) {
	if err != nil {
		return "", err
	}

	// 定义响应结构体用于解析JSON
	var res httpx.Response
	if err := json.Unmarshal(v, &res); err != nil {
		return "", err
	}

	if res.Code != 200 {
		return "", errors.New(res.Msg)
	}

	switch chatType {
	case domain.TodoAdd:
		return Success, err
	}

	data := domain.ChatResp{
		ChatType: int(chatType),
		Data:     res.Data,
	}

	// 将响应结构体序列化为JSON字符串
	d, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	return SuccessWithData + string(d) + "\n\n\n", nil
}
