package toolx

import (
	"ai/internal/logic/chatinternal/toolx/approval"
	"ai/internal/model"
	"ai/internal/svc"
	"ai/pkg/langchain/outputparserx"
	"context"

	"github.com/tmc/langchaingo/callbacks"
)

type ApprovalAdd struct {
	svc          *svc.ServiceContext
	Callback     callbacks.Handler
	outputparser outputparserx.Structured
}

// NewApprovalAdd 创建一个新的审批添加处理器实例
func NewApprovalAdd(svc *svc.ServiceContext) *ApprovalAdd {
	return &ApprovalAdd{
		svc:      svc,
		Callback: svc.Callbacks,
		outputparser: outputparserx.NewStructured([]outputparserx.ResponseSchema{
			{
				Name: "type",
				Description: "approval type; enum : 1. General, 2. Leave Approval, 3. " +
					"Card replacement Approval, 4. Go out Approval",
				Type: "int",
			}, {
				Name:        "input",
				Description: "The user's original input",
			},
		}),
	}
}

func (a *ApprovalAdd) Name() string {
	return "approval_add"
}

func (a *ApprovalAdd) Description() string {
	return `
	a approval add interface
	use when you need to create a approval.
	keep Chinese output.
` + a.outputparser.GetFormatInstructions()
}

// Call 执行审批创建操作的核心方法
func (a *ApprovalAdd) Call(ctx context.Context, input string) (string, error) {
	if a.Callback != nil {
		a.Callback.HandleText(ctx, "approval add start input : "+input)
	}

	// 解析用户输入，按照预设的结构化格式提取数据
	out, err := a.outputparser.Parse(input)
	if err != nil {
		return "", err
	}
	data := out.(map[string]any)

	// 提取审批类型（默认为0，后续可能需要处理默认值）
	var approvalType float64
	if t, ok := data["type"]; ok {
		approvalType = t.(float64)
	}
	// 提取用户输入内容（如果存在）
	if v, ok := data["input"]; ok {
		input = v.(string)
	}

	// 设置公司工作时间信息（这里覆盖了原有input，可能需要根据实际业务调整）
	input = "The company's working hours are normal working hours of 8 hours a day and 40 hours a week; Monday to Friday 9:30-11:30 13:00-18:00"

	// 根据审批类型创建对应的审批实例
	ap, err := approval.NewApproval(a.svc, model.ApprovalType(approvalType))
	if err != nil {
		return "", err
	}

	// 创建审批记录并获取ID
	id, err := ap.Create(ctx, input)
	if err != nil {
		return "", err
	}

	// 返回成功信息及创建的审批ID
	return Success + "\n created approval id : " + id, nil
}
