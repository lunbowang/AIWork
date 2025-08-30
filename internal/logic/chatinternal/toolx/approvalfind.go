package toolx

import (
	"ai/internal/domain"
	"ai/internal/svc"
	"ai/pkg/curl"
	"ai/pkg/langchain/outputparserx"
	"ai/token"
	"context"

	"github.com/tmc/langchaingo/callbacks"
)

type ApprovalFind struct {
	svc          *svc.ServiceContext
	Callback     callbacks.Handler
	outputparser outputparserx.Structured
}

func NewApprovalFind(svc *svc.ServiceContext) *ApprovalFind {
	return &ApprovalFind{
		svc:      svc,
		Callback: svc.Callbacks,
		outputparser: outputparserx.NewStructured([]outputparserx.ResponseSchema{
			{
				Name:        "type",
				Description: "approval type; enum : 0. General, 1. Leave, 2. Card replacement, 3. Go out, 4. Reimbursement, 5. Payment, 6. Purchase, 7. Collection; number to be completed",
				Type:        "int",
			}, {
				Name:        "id",
				Description: "approval id",
			}, {
				Name:        "status",
				Description: "approval status; enum : 0. No beginning, 1. In progress 2. Done-Passed, 3. Revocation, 4. refused; number to be completed",
			}, {
				Name:        "createId",
				Description: "id of creator",
			},
		}),
	}
}

func (a *ApprovalFind) Name() string {
	return "approval_find"
}

func (a *ApprovalFind) Description() string {
	return `
	a approval find interface.
	use when you need to find a approval.
	If the condition is null, return {}
 	keep Chinese output.` + a.outputparser.GetFormatInstructions()
}

func (a *ApprovalFind) Call(ctx context.Context, input string) (string, error) {
	if a.Callback != nil {
		a.Callback.HandleText(ctx, "approval find start input : "+input)
	}

	out, err := a.outputparser.Parse(input)
	if err != nil {
		return "", err
	}

	data := out.(map[string]any)
	if data == nil {
		data = make(map[string]any)
	}
	if data["createId"] == nil {
		data["createId"] = token.GetUId(ctx)
	}

	res, err := curl.GetRequest(token.GetTokenStr(ctx), a.svc.Config.Host+"/v1/approval/list", data)
	if err != nil {
		return "", err
	}

	if a.Callback != nil {
		a.Callback.HandleText(ctx, "approval find end data : "+string(res))
	}

	return ResParser(res, domain.ApprovalFind, nil)
}
