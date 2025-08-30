package api

import (
	"github.com/gin-gonic/gin"

	"ai/internal/domain"
	"ai/internal/logic"
	"ai/internal/svc"
	"ai/pkg/httpx"
)

type Approval struct {
	svcCtx   *svc.ServiceContext
	approval logic.Approval
}

func NewApproval(svcCtx *svc.ServiceContext, approval logic.Approval) *Approval {
	return &Approval{
		svcCtx:   svcCtx,
		approval: approval,
	}
}

func (h *Approval) InitRegister(engine *gin.Engine) {
	g := engine.Group("v1/approval", h.svcCtx.Jwt.Handler)
	g.GET("/:id", h.Info)
	g.POST("", h.Create)
	g.PUT("/dispose", h.Dispose)
	g.GET("/list", h.List)
}

func (h *Approval) Info(ctx *gin.Context) {
	var req domain.IdPathReq
	if err := httpx.BindAndValidate(ctx, &req); err != nil {
		httpx.FailWithErr(ctx, err)
		return
	}

	res, err := h.approval.Info(ctx.Request.Context(), &req)
	if err != nil {
		httpx.FailWithErr(ctx, err)
	} else {
		httpx.OkWithData(ctx, res)
	}
}

func (h *Approval) Create(ctx *gin.Context) {
	var req domain.Approval
	if err := httpx.BindAndValidate(ctx, &req); err != nil {
		httpx.FailWithErr(ctx, err)
		return
	}

	res, err := h.approval.Create(ctx.Request.Context(), &req)
	if err != nil {
		httpx.FailWithErr(ctx, err)
	} else {
		httpx.OkWithData(ctx, res)
	}
}

func (h *Approval) Dispose(ctx *gin.Context) {
	var req domain.DisposeReq
	if err := httpx.BindAndValidate(ctx, &req); err != nil {
		httpx.FailWithErr(ctx, err)
		return
	}

	err := h.approval.Dispose(ctx.Request.Context(), &req)
	if err != nil {
		httpx.FailWithErr(ctx, err)
	} else {
		httpx.Ok(ctx)
	}
}

func (h *Approval) List(ctx *gin.Context) {
	var req domain.ApprovalListReq
	if err := httpx.BindAndValidate(ctx, &req); err != nil {
		httpx.FailWithErr(ctx, err)
		return
	}

	res, err := h.approval.List(ctx.Request.Context(), &req)
	if err != nil {
		httpx.FailWithErr(ctx, err)
	} else {
		httpx.OkWithData(ctx, res)
	}
}
