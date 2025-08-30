package chatinternal

import (
	"ai/internal/domain"
	"ai/internal/model"
	"ai/internal/svc"
	"ai/pkg/langchain"
	"ai/pkg/langchain/outputparserx"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/tmc/langchaingo/prompts"

	"github.com/tmc/langchaingo/chains"
)

// Task 聊天总结结果的任务结构体
type Task struct {
	Type    int
	Title   string
	Content string
}

// ChatLogHandle 聊天记录处理核心结构体
type ChatLogHandle struct {
	svc    *svc.ServiceContext
	chains chains.Chain
	out    outputparserx.Structured
}

// NewChatLogHandle 创建ChatLogHandle实例
func NewChatLogHandle(svc *svc.ServiceContext) *ChatLogHandle {
	return &ChatLogHandle{
		svc:    svc,
		chains: chains.NewLLMChain(svc.LLMs, prompts.NewPromptTemplate(_defaultChatLogPrompts, []string{"input"})),
		out:    outputparserx.Structured{},
	}
}

func (c *ChatLogHandle) Name() string {
	return "chat_log"
}

func (c *ChatLogHandle) Description() string {
	return "used to summarize and analyze the content of a chat session"

}

func (c *ChatLogHandle) Chains() chains.Chain {
	return chains.NewTransform(c.transform, nil, nil)
}

// transform 核心转换方法：实现“聊天记录→结构化总结结果”的完整流程
func (c *ChatLogHandle) transform(ctx context.Context, inputs map[string]any,
	opts ...chains.ChainCallOption) (map[string]any, error) {
	// 1. 从输入参数中提取“会话ID（relationId）”：必传参数，无则返回错误
	var cid string
	if id, ok := inputs["relationId"].(string); !ok {
		return nil, errors.New("请确定需要总结的会话对象")
	} else {
		cid = id
	}

	// 2. 处理时间范围：若输入未传startTime/endTime，默认查询“近24小时”的记录
	startTime, endTime := setTimeRange(inputs)

	// 3. 查询指定会话、指定时间范围内的聊天记录（调用chatLog方法）
	msgs, err := c.chatLog(ctx, cid, startTime, endTime)
	if err != nil {
		return nil, err
	}

	// 4. 调用LLM链总结聊天记录：将聊天记录作为input传给大模型
	res, err := chains.Call(ctx, c.chains, map[string]any{
		"input": msgs,
	})
	if err != nil {
		return nil, err
	}

	// 5. 提取LLM输出结果：需确保输出是string类型，否则返回“无效输出”错误
	text, ok := res[langchain.OutPut].(string)
	if !ok {
		return nil, chains.ErrInvalidOutputValues
	}

	// 6. 将LLM输出的JSON字符串解析为Task数组（核心：转换为业务结构体）
	var data []*Task
	if err := json.Unmarshal([]byte(text), &data); err != nil {
		return nil, err
	}

	// 7. 构造最终业务响应：封装为domain.ChatResp格式（包含聊天类型+任务数据）
	b, err := json.Marshal(domain.ChatResp{
		ChatType: domain.ChatLog,
		Data:     data,
	})
	if err != nil {
		return nil, err
	}

	// 8. 返回标准格式的输出结果：key固定为langchain.OutPut，value为JSON字符串
	return map[string]any{
		langchain.OutPut: string(b),
	}, nil
}

// chatLog 查询指定会话、时间范围内的聊天记录，并拼接为“姓名(ID): 消息内容”的文本格式
func (c *ChatLogHandle) chatLog(ctx context.Context, cid string, startTime, endTime int64) (string, error) {
	// 1. 调用ChatlogModel查询聊天记录（按发送时间筛选）
	list, err := c.svc.ChatlogModel.ListBySendTime(ctx, cid, startTime, endTime)
	if err != nil {
		return "", err
	}

	// 2. 定义聊天记录拼接模板：姓名(ID): 消息内容 + 换行
	chatStr := "%s(%s): %s\n"

	// 3. 初始化字符串构建器（高效拼接大量文本，避免string拼接的内存浪费）
	var (
		res    strings.Builder
		record = make(map[string]*model.User)
	)

	// 4. 遍历聊天记录，拼接文本并缓存用户信息
	for i, _ := range list {
		var u *model.User
		if v, ok := record[list[i].SendId]; ok {
			u = v
		} else {
			t, err := c.svc.UserModel.FindOne(ctx, list[i].SendId)
			if err != nil {
				return "", err
			}

			u = t
			record[list[i].SendId] = t
		}

		// 5. 按模板拼接当前聊天记录，写入字符串构建器
		res.Write([]byte(fmt.Sprintf(chatStr, u.Name, u.ID.Hex(), list[i].MsgContent)))
	}

	// 6. 返回拼接完成的聊天记录文本
	return res.String(), nil
}

// setTimeRange 处理输入的时间范围参数，提供默认值
func setTimeRange(input map[string]any) (int64, int64) {
	var startTime, endTime int64

	// 获取当前时间戳（单位秒）：作为默认时间范围的基准
	cuurentTime := time.Now().Unix()

	// 处理startTime：若输入有合法的startTime（int类型且>0），则使用；否则默认“24小时前”
	if start, ok := input["startTime"].(int); ok && start > 0 {
		startTime = int64(start)
	} else {
		// 获取前一天的消息
		startTime = cuurentTime - 24*3600
	}

	// 处理endTime：若输入有合法的endTime（int类型且>0），则使用；否则默认“当前时间”
	if end, ok := input["endTime"].(int); ok && end > 0 {
		endTime = int64(end)
	} else {
		// 获取前一天的消息
		endTime = cuurentTime
	}

	return startTime, endTime
}
