package api

import (
	"github.com/gin-gonic/gin"

	"ai/internal/domain"
	"ai/internal/logic"
	"ai/internal/svc"
	"ai/pkg/httpx"
)

type Department struct {
	svcCtx     *svc.ServiceContext
	department logic.Department
}

func NewDepartment(svcCtx *svc.ServiceContext, department logic.Department) *Department {
	return &Department{
		svcCtx:     svcCtx,
		department: department,
	}
}

func (h *Department) InitRegister(engine *gin.Engine) {
	g := engine.Group("v1/dep", h.svcCtx.Jwt.Handler)
	g.GET("/soa", h.Soa)
	g.GET("/:id", h.Info)
	g.POST("", h.Create)
	g.PUT("", h.Edit)
	g.DELETE("/:id", h.Delete)
	g.POST("/user", h.SetDepUsers)
	g.GET("/user/:id", h.DepUserInfo)
}

func (h *Department) Soa(ctx *gin.Context) {
	res, err := h.department.Soa(ctx.Request.Context())
	if err != nil {
		httpx.FailWithErr(ctx, err)
	} else {
		httpx.OkWithData(ctx, res)
	}
}

func (h *Department) Info(ctx *gin.Context) {
	var req domain.IdPathReq
	if err := httpx.BindAndValidate(ctx, &req); err != nil {
		httpx.FailWithErr(ctx, err)
		return
	}

	res, err := h.department.Info(ctx.Request.Context(), &req)
	if err != nil {
		httpx.FailWithErr(ctx, err)
	} else {
		httpx.OkWithData(ctx, res)
	}
}

func (h *Department) Create(ctx *gin.Context) {
	var req domain.Department
	if err := httpx.BindAndValidate(ctx, &req); err != nil {
		httpx.FailWithErr(ctx, err)
		return
	}

	err := h.department.Create(ctx.Request.Context(), &req)
	if err != nil {
		httpx.FailWithErr(ctx, err)
	} else {
		httpx.Ok(ctx)
	}
}

func (h *Department) Edit(ctx *gin.Context) {
	var req domain.Department
	if err := httpx.BindAndValidate(ctx, &req); err != nil {
		httpx.FailWithErr(ctx, err)
		return
	}

	err := h.department.Edit(ctx.Request.Context(), &req)
	if err != nil {
		httpx.FailWithErr(ctx, err)
	} else {
		httpx.Ok(ctx)
	}
}

func (h *Department) Delete(ctx *gin.Context) {
	var req domain.IdPathReq
	if err := httpx.BindAndValidate(ctx, &req); err != nil {
		httpx.FailWithErr(ctx, err)
		return
	}

	err := h.department.Delete(ctx.Request.Context(), &req)
	if err != nil {
		httpx.FailWithErr(ctx, err)
	} else {
		httpx.Ok(ctx)
	}
}

func (h *Department) SetDepUsers(ctx *gin.Context) {
	var req domain.SetDepUser
	if err := httpx.BindAndValidate(ctx, &req); err != nil {
		httpx.FailWithErr(ctx, err)
		return
	}

	err := h.department.SetDepUsers(ctx.Request.Context(), &req)
	if err != nil {
		httpx.FailWithErr(ctx, err)
	} else {
		httpx.Ok(ctx)
	}
}

// 获取用户的部门信息
func (h *Department) DepUserInfo(ctx *gin.Context) {
	var req domain.IdPathReq
	if err := httpx.BindAndValidate(ctx, &req); err != nil {
		httpx.FailWithErr(ctx, err)
		return
	}

	res, err := h.department.DepUserInfo(ctx.Request.Context(), &req)
	if err != nil {
		httpx.FailWithErr(ctx, err)
	} else {
		httpx.OkWithData(ctx, res)
	}
}
