package logic

import (
	"ai/internal/logic/chatinternal"
	"ai/internal/model"
	"ai/pkg/langchain"
	"ai/pkg/langchain/memoryx"
	"ai/pkg/langchain/router"
	"ai/pkg/langchain/voice"
	"ai/token"
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"regexp"
	"sort"
	"time"

	"github.com/tmc/langchaingo/schema"

	"gitee.com/dn-jinmin/tlog"

	"github.com/tmc/langchaingo/chains"

	"ai/internal/domain"
	"ai/internal/svc"
)

type Chat interface {
	PrivateChat(ctx context.Context, req *domain.Message) error
	GroupChat(ctx context.Context, req *domain.Message) (uids []string, err error)
	AIChat(ctx context.Context, req *domain.ChatReq) (resp *domain.ChatResp, err error)
	File(ctx context.Context, req []*domain.FileResp) (err error)
}

type chat struct {
	svc    *svc.ServiceContext
	memory schema.Memory
	router *router.Router
	voice  *voice.Voice
}

func NewChat(svc *svc.ServiceContext) Chat {
	handlers := []router.Handler{
		chatinternal.NewTodoHandle(svc),
		chatinternal.NewKnowledge(svc),
		chatinternal.NewApprovalHandle(svc),
		chatinternal.NewChatLogHandle(svc),
	}

	memory := memoryx.NewMemoryx(func() schema.Memory {
		m := memoryx.NewSummaryBuffer(svc.LLMs, 50, memoryx.WithCallback(svc.Callbacks),
			memoryx.WithOutParser(memoryOutput))

		m.InputKey = langchain.Input
		return m
	})
	return &chat{
		svc:    svc,
		memory: memory,
		voice:  voice.NewVoice(svc.OpenaiClient),
		router: router.NewRouter(svc.LLMs, handlers,
			router.WithMemory(memory),
			router.Withcallback(svc.Callbacks),
			router.WithEmptyHandler(chatinternal.NewDefaultHandler(svc)),
		),
	}
}

func (l *chat) PrivateChat(ctx context.Context, req *domain.Message) error {
	return l.chatlog(ctx, req)
}

func (l *chat) GroupChat(ctx context.Context, req *domain.Message) (uids []string, err error) {
	req.ConversationId = "all"
	if err := l.chatlog(ctx, req); err != nil {
		return nil, err
	}

	return nil, err
}

func (l *chat) chatlog(ctx context.Context, req *domain.Message) error {
	sendId := req.SendId

	chatlog := model.Chatlog{
		ConversationId: req.ConversationId,
		SendId:         sendId,
		RecvId:         req.RecvId,
		ChatType:       model.ChatType(req.ChatType),
		MsgContent:     req.Content,
		SendTime:       time.Now().Unix(),
	}

	if chatlog.ConversationId == "" {
		chatlog.ConversationId = GenerateUniqueID(sendId, req.RecvId)
	}

	return l.svc.ChatlogModel.Insert(ctx, &chatlog)
}

func (l *chat) AIChat(ctx context.Context, req *domain.ChatReq) (resp *domain.ChatResp, err error) {
	uid := token.GetUId(ctx)
	ctx = context.WithValue(ctx, langchain.ChatId, uid)

	if req.ChatType > 0 {
		return l.basicService(ctx, req)
	}

	if len(req.Prompts) == 0 && len(req.Voice) > 0 {
		prompts, err := l.voice.Transcriptions(ctx, req.Voice, "")
		if err != nil {
			return nil, err
		}

		req.Prompts = prompts
	}

	return l.aiService(ctx, req)
}

func (l *chat) aiService(ctx context.Context, req *domain.ChatReq) (resp *domain.ChatResp, err error) {
	v, err := chains.Call(ctx, l.router, map[string]any{
		langchain.Input: req.Prompts,
		"relationId":    req.RelationId,
		"startTime":     req.StartTime,
		"endTime":       req.EndTime,
	}, chains.WithCallback(l.svc.Callbacks))
	if err != nil {
		fmt.Println("\n aiService chains ---- err : ", err.Error())
		var res domain.ChatResp
		re := regexp.MustCompile(`{.*}`)
		matches := re.FindString(err.Error())
		matches = regexp.MustCompile(`\\`).ReplaceAllString(matches, "")
		if e := json.Unmarshal([]byte(matches), &res); e != nil {
			return nil, e
		}
		return &res, nil
	}
	fmt.Println("\n aiService chains ---- thoroughAnalysis ")
	return l.thoroughAnalysis(ctx, v)
}

func (l *chat) thoroughAnalysis(ctx context.Context, out map[string]any) (*domain.ChatResp, error) {
	var str string
	if _, ok := out[langchain.OutPut]; !ok {
		return &domain.ChatResp{
			ChatType: domain.DefaultHandler,
			Data:     str,
		}, nil
	}
	// 初步解析
	str = out[langchain.OutPut].(string)
	var res domain.ChatResp
	if err := json.Unmarshal([]byte(str), &res); err == nil {
		return &res, nil
	}
	// 再增加增加的方式验证
	re := regexp.MustCompile(`{.*}`)
	matches := re.FindString(str)
	matches = regexp.MustCompile(`\\`).ReplaceAllString(matches, "")
	if err := json.Unmarshal([]byte(matches), &res); err == nil {
		return &res, nil
	}

	// 还可以通过于最终方案 大模型 处理 或者 直接输出
	return &domain.ChatResp{
		ChatType: domain.DefaultHandler,
		Data:     out,
	}, nil

}

func (l *chat) basicService(ctx context.Context, req *domain.ChatReq) (resp *domain.ChatResp, err error) {
	return nil, err
}

func memoryOutput(ctx context.Context, v string) string {
	var res domain.ChatResp
	if err := json.Unmarshal([]byte(v), &res); err != nil {
		tlog.ErrorfCtx(ctx, "memoryOutput", "v %s, err %s", v, err.Error())
		return v
	}

	tlog.InfoCtx(ctx, "memoryOutput", v)

	switch res.Data.(type) {
	case string:
		return res.Data.(string)
	default:
		return v
	}
}

func (l *chat) File(ctx context.Context, files []*domain.FileResp) (err error) {
	uid := token.GetUId(ctx)
	ctx = context.WithValue(ctx, langchain.ChatId, uid)

	data := make([]*domain.ChatFile, 0, len(files))
	for _, file := range files {
		data = append(data, &domain.ChatFile{
			Path: file.File,
			Name: file.Filename,
			Time: time.Now(),
		})
	}

	b, err := json.Marshal(data)
	if err != nil {
		return err
	}

	err = l.memory.SaveContext(ctx, map[string]any{
		langchain.Input: string(b),
	}, map[string]any{
		langchain.Input: "uploaded files",
	})

	// test
	memoryContent, err := l.memory.LoadMemoryVariables(ctx, map[string]any{})
	if err != nil {
		return err
	}
	fmt.Println(memoryContent)

	return
}

// GenerateUniqueID 根据传递的两个字符串 ID 生成唯一的 ID
func GenerateUniqueID(id1, id2 string) string {
	// 将两个 ID 放入切片中
	ids := []string{id1, id2}

	// 对 IDs 切片进行排序
	sort.Strings(ids)

	// 将排序后的 ID 组合起来
	combined := ids[0] + ids[1]

	// 创建 SHA-256 哈希对象
	hasher := sha256.New()

	// 写入合并后的字符串
	hasher.Write([]byte(combined))

	// 计算哈希值
	hash := hasher.Sum(nil)

	// 返回哈希值的十六进制字符串表示
	return base64.RawStdEncoding.EncodeToString(hash)[:22] // 可以选择更短的长度
}
