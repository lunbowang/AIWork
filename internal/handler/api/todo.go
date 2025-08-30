package api

import (
	"github.com/gin-gonic/gin"

	"ai/internal/domain"
	"ai/internal/logic"
	"ai/internal/svc"
	"ai/pkg/httpx"
)

type Todo struct {
	svcCtx *svc.ServiceContext
	todo   logic.Todo
}

func NewTodo(svcCtx *svc.ServiceContext, todo logic.Todo) *Todo {
	return &Todo{
		svcCtx: svcCtx,
		todo:   todo,
	}
}

func (h *Todo) InitRegister(engine *gin.Engine) {
	g := engine.Group("v1/todo", h.svcCtx.Jwt.Handler)
	g.GET("/:id", h.Info)
	g.POST("", h.Create)
	g.PUT("", h.Edit)
	g.DELETE("/:id", h.Delete)
	g.POST("/finish", h.Finish)
	g.POST("/record", h.CreateRecord)
	g.GET("/list", h.List)
}

func (h *Todo) Info(ctx *gin.Context) {
	var req domain.IdPathReq
	if err := httpx.BindAndValidate(ctx, &req); err != nil {
		httpx.FailWithErr(ctx, err)
		return
	}

	res, err := h.todo.Info(ctx.Request.Context(), &req)
	if err != nil {
		httpx.FailWithErr(ctx, err)
	} else {
		httpx.OkWithData(ctx, res)
	}
}

func (h *Todo) Create(ctx *gin.Context) {
	var req domain.Todo
	if err := httpx.BindAndValidate(ctx, &req); err != nil {
		httpx.FailWithErr(ctx, err)
		return
	}

	res, err := h.todo.Create(ctx.Request.Context(), &req)
	if err != nil {
		httpx.FailWithErr(ctx, err)
	} else {
		httpx.OkWithData(ctx, res)
	}
}

func (h *Todo) Edit(ctx *gin.Context) {
	var req domain.Todo
	if err := httpx.BindAndValidate(ctx, &req); err != nil {
		httpx.FailWithErr(ctx, err)
		return
	}

	err := h.todo.Edit(ctx.Request.Context(), &req)
	if err != nil {
		httpx.FailWithErr(ctx, err)
	} else {
		httpx.Ok(ctx)
	}
}

func (h *Todo) Delete(ctx *gin.Context) {
	var req domain.IdPathReq
	if err := httpx.BindAndValidate(ctx, &req); err != nil {
		httpx.FailWithErr(ctx, err)
		return
	}

	err := h.todo.Delete(ctx.Request.Context(), &req)
	if err != nil {
		httpx.FailWithErr(ctx, err)
	} else {
		httpx.Ok(ctx)
	}
}

// Finish 完成待办
func (h *Todo) Finish(ctx *gin.Context) {
	var req domain.FinishedTodoReq
	if err := httpx.BindAndValidate(ctx, &req); err != nil {
		httpx.FailWithErr(ctx, err)
		return
	}

	err := h.todo.Finish(ctx.Request.Context(), &req)
	if err != nil {
		httpx.FailWithErr(ctx, err)
	} else {
		httpx.Ok(ctx)
	}
}

// CreateRecord 用户操作记录
func (h *Todo) CreateRecord(ctx *gin.Context) {
	var req domain.TodoRecord
	if err := httpx.BindAndValidate(ctx, &req); err != nil {
		httpx.FailWithErr(ctx, err)
		return
	}

	err := h.todo.CreateRecord(ctx.Request.Context(), &req)
	if err != nil {
		httpx.FailWithErr(ctx, err)
	} else {
		httpx.Ok(ctx)
	}
}

// 待办列表
func (h *Todo) List(ctx *gin.Context) {
	var req domain.TodoListReq
	if err := httpx.BindAndValidate(ctx, &req); err != nil {
		httpx.FailWithErr(ctx, err)
		return
	}

	res, err := h.todo.List(ctx.Request.Context(), &req)
	if err != nil {
		httpx.FailWithErr(ctx, err)
	} else {
		httpx.OkWithData(ctx, res)
	}
}
