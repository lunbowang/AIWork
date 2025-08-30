package approval

import (
	"ai/internal/domain"
	"ai/internal/model"
	"ai/internal/svc"
	"ai/pkg/curl"
	"ai/pkg/langchain"
	"ai/pkg/langchain/outputparserx"
	"ai/pkg/xerr"
	"ai/token"
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/go-viper/mapstructure/v2"

	"github.com/tmc/langchaingo/prompts"

	"github.com/tmc/langchaingo/chains"
)

// Approval 审批接口，定义了所有审批类型都必须实现的创建方法
type Approval interface {
	Create(ctx context.Context, input string) (string, error)
}

// Approvals 审批类型与对应审批实例的映射表
// 用于根据审批类型快速获取对应的审批处理实例
var Approvals map[model.ApprovalType]Approval

// NewApproval 根据审批类型创建对应的审批实例
func NewApproval(svc *svc.ServiceContext, approvalType model.ApprovalType) (Approval, error) {
	if Approvals == nil {
		Approvals = map[model.ApprovalType]Approval{
			model.LeaveApproval:    NewLeave(svc),
			model.MakeCardApproval: NewMakeCard(svc),
			model.GoOutApproval:    NewGoOut(svc),
		}
	}

	a := Approvals[approvalType]
	if a == nil {
		return nil, errors.New("不存在该审批类型" + fmt.Sprintf("%v", approvalType))
	}

	return a, nil
}

type MakeCard struct {
	svc          *svc.ServiceContext
	c            chains.Chain
	outputparser outputparserx.Structured
}

func NewMakeCard(svc *svc.ServiceContext) *MakeCard {
	output := outputparserx.NewStructured([]outputparserx.ResponseSchema{
		{
			Name:        "date",
			Description: "filling time,data application time stamp,  such as 1720921573",
			Type:        "int64",
		}, {
			Name:        "reason",
			Description: "reason for replacement card",
		}, {
			Name:        "day",
			Description: "replacement date, such as 20221011",
			Type:        "int64",
		}, {
			Name:        "workCheckType",
			Description: "replacement card type; enum : 1. Work card. 2. Work card; number to be completed",
			Type:        "int",
		},
	})
	return &MakeCard{
		svc: svc,
		c: chains.NewLLMChain(svc.LLMs, prompts.NewPromptTemplate(
			_defaultCreateApprovalTemplate+output.GetFormatInstructions(), []string{"input"},
		)),
		outputparser: output,
	}
}

// Create 实现Approval接口，创建补卡审批
func (m *MakeCard) Create(ctx context.Context, input string) (string, error) {
	// 调用LLM链处理输入，获取自然语言处理结果
	out, err := chains.Predict(ctx, m.c, map[string]any{
		langchain.Input: input,
	})
	if err != nil {
		return "", err
	}

	// 将LLM输出解析为结构化数据
	v, err := m.outputparser.Parse(out)
	if err != nil {
		return "", err
	}

	// 将结构化数据映射到补卡审批领域模型
	var data domain.MakeCard
	if err := mapstructure.Decode(v, &data); err != nil {
		return "", err
	}

	// 构建审批请求对象
	req := domain.Approval{
		Type:     int(model.MakeCardApproval),
		MakeCard: &data,
	}
	addRes, err := curl.PostRequest(token.GetTokenStr(ctx), m.svc.Config.Host+"/v1/approval", req)
	if err != nil {
		return "", err
	}

	var idResp domain.IdRespInfo
	if err := json.Unmarshal(addRes, &idResp); err != nil {
		return "", xerr.WithMessage(err, "")
	}

	return idResp.Data.Id, err

}

type Leave struct {
	svc          *svc.ServiceContext
	c            chains.Chain
	outPutParser outputparserx.Structured
}

func NewLeave(svc *svc.ServiceContext) *Leave {
	output := outputparserx.NewStructured([]outputparserx.ResponseSchema{
		{
			Name:        "type",
			Description: "type of leave; enum 0. Personal leave, 1. Vacation, 2. Sick leave, 3. Annual leave, 4. Maternity leave, 5. Paternity leave, 6. Marriage leave, 7. Bereavement leave, 8. Breastfeeding leave; number to be completed",
			Type:        "int",
		}, {
			Name:        "startTime",
			Description: "leave start time,data application time stamp,  such as 1720921573",
			Type:        "int64",
		}, {
			Name:        "startTime",
			Description: "leave end time,data application time stamp,  such as 1720921573",
			Type:        "int64",
		}, {
			Name:        "reason",
			Description: "Reason for leave",
		}, {
			Name:        "timeType",
			Description: "Leave time type; enum 1. Hours, 2. Days; Use the day type for more than 24 hours, and use the hour type for less than 23 hours; number to be completed",
			Type:        "int64",
		},
	})
	return &Leave{
		svc: svc,
		c: chains.NewLLMChain(svc.LLMs, prompts.NewPromptTemplate(
			_defaultCreateApprovalTemplate+output.GetFormatInstructions(), []string{"input"},
		)),
	}
}

func (m *Leave) Create(ctx context.Context, input string) (string, error) {
	out, err := chains.Predict(ctx, m.c, map[string]any{
		langchain.Input: input,
	}, chains.WithCallback(m.svc.Callbacks))
	if err != nil {
		return "", xerr.WithMessage(err, "chains.Predict : "+input)
	}

	v, err := m.outPutParser.Parse(out)
	if err != nil {
		return "", xerr.WithMessage(err, "m.outPutParser.Parse")
	}

	var data domain.Leave
	if err := mapstructure.Decode(v, &data); err != nil {
		return "", xerr.WithMessage(err, "domain.GoOut")
	}

	req := domain.Approval{
		Type:  int(model.LeaveApproval),
		Leave: &data,
	}

	fmt.Println("提交请假审批 ： ", req, " \n ", req.Leave)

	addRes, err := curl.PostRequest(token.GetTokenStr(ctx), m.svc.Config.Host+"/v1/approval", req)
	var idResp domain.IdRespInfo
	if err := json.Unmarshal(addRes, &idResp); err != nil {
		return "", xerr.WithMessage(err, "")
	}

	return idResp.Data.Id, err
}

type GoOut struct {
	svc          *svc.ServiceContext
	c            chains.Chain
	outPutParser outputparserx.Structured
}

func NewGoOut(svc *svc.ServiceContext) *MakeCard {
	output := outputparserx.NewStructured([]outputparserx.ResponseSchema{
		{
			Name:        "startTime",
			Description: "go out start time,data application time stamp, such as 1720921573",
			Type:        "int64",
		}, {
			Name:        "startTime",
			Description: "go out end time,data application time stamp, such as 1720921573",
			Type:        "int64",
		}, {
			Name:        "reason",
			Description: "Reason for go out",
		},
	})
	return &MakeCard{
		svc: svc,
		c: chains.NewLLMChain(svc.LLMs, prompts.NewPromptTemplate(
			_defaultCreateApprovalTemplate+output.GetFormatInstructions(), []string{"input"},
		)),
	}
}

func (m *GoOut) Create(ctx context.Context, input string) (string, error) {

	out, err := chains.Predict(ctx, m.c, map[string]any{
		langchain.Input: input,
	})
	if err != nil {
		return "", xerr.WithMessage(err, "chains.Predict : "+input)
	}

	v, err := m.outPutParser.Parse(out)
	if err != nil {
		return "", xerr.WithMessage(err, " m.outPutParser.Parse")
	}

	var data domain.GoOut
	if err := mapstructure.Decode(v, &data); err != nil {
		return "", xerr.WithMessage(err, "domain.GoOut")
	}

	req := domain.Approval{
		Type:  int(model.GoOutApproval),
		GoOut: &data,
	}

	addRes, err := curl.PostRequest(token.GetTokenStr(ctx), m.svc.Config.Host+"/v1/approval", req)
	var idResp domain.IdRespInfo
	if err := json.Unmarshal(addRes, &idResp); err != nil {
		return "", xerr.WithMessage(err, "")
	}

	return idResp.Data.Id, err
}
