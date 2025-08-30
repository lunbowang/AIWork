package api

import (
	"ai/internal/logic"
	"ai/internal/svc"
)

func initHandler(svc *svc.ServiceContext) []Handler {
	// new logics
	var (
		departmentLogic = logic.NewDepartment(svc)
		todoLogic       = logic.NewTodo(svc)
		approvalLogic   = logic.NewApproval(svc)
		chatLogic       = logic.NewChat(svc)
		userLogic       = logic.NewUser(svc)
	)

	// new handlers
	var (
		todo       = NewTodo(svc, todoLogic)
		approval   = NewApproval(svc, approvalLogic)
		chat       = NewChat(svc, chatLogic)
		upload     = NewUpload(svc, chatLogic)
		user       = NewUser(svc, userLogic)
		department = NewDepartment(svc, departmentLogic)
	)

	return []Handler{
		todo,
		approval,
		chat,
		upload,
		user,
		department,
	}
}
