package chatinternal

import (
	"ai/internal/logic/chatinternal/toolx"
	"ai/internal/svc"

	"github.com/tmc/langchaingo/tools"
)

type ApprovalHandle struct {
	*baseChat
}

func NewApprovalHandle(svc *svc.ServiceContext) *ApprovalHandle {
	return &ApprovalHandle{
		baseChat: NewBaseChat(svc, []tools.Tool{
			toolx.NewApprovalAdd(svc),
			toolx.NewApprovalFind(svc),
		}),
	}
}
func (t *ApprovalHandle) Name() string {
	return "approval"
}

func (t *ApprovalHandle) Description() string {
	return "This is about approval matters. Such as sick leave, personal leave, going out, etc.\n\n"
}
