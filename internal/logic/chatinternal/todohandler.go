package chatinternal

import (
	"ai/internal/logic/chatinternal/toolx"
	"ai/internal/svc"

	"github.com/tmc/langchaingo/tools"
)

type TodoHandle struct {
	*baseChat
}

func NewTodoHandle(svc *svc.ServiceContext) *TodoHandle {
	return &TodoHandle{
		baseChat: NewBaseChat(svc, []tools.Tool{
			toolx.NewTodoAdd(svc),
			toolx.NewTodoFind(svc),
		}),
	}
}

func (t *TodoHandle) Name() string {
	return "todo"
}

func (t *TodoHandle) Description() string {
	return "suitable for todo processing, such as todo creation, query, modification, dele tion, etc"
}
