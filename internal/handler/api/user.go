package api

import (
	"ai/internal/domain"
	"ai/internal/logic"
	"ai/internal/svc"
	"ai/pkg/httpx"

	"github.com/gin-gonic/gin"
)

type User struct {
	svcCtx *svc.ServiceContext
	user   logic.User
}

func NewUser(svcCtx *svc.ServiceContext, user logic.User) *User {
	return &User{
		svcCtx: svcCtx,
		user:   user,
	}
}

func (h *User) InitRegister(engine *gin.Engine) {
	g0 := engine.Group("v1/user")
	g0.POST("/login", h.Login)

	g1 := engine.Group("/v1/user", h.svcCtx.Jwt.Handler)
	g1.GET("/:id", h.Info)
	g1.POST("", h.Create)
	g1.PUT("", h.Edit)
	g1.DELETE("/:id", h.Delete)
	g1.GET("/list", h.List)
	g1.POST("/password", h.UpPassword)
}

func (h *User) Login(ctx *gin.Context) {
	var req domain.LoginReq
	if err := httpx.BindAndValidate(ctx, &req); err != nil {
		httpx.FailWithErr(ctx, err)
		return
	}

	res, err := h.user.Login(ctx.Request.Context(), &req)
	if err != nil {
		httpx.FailWithErr(ctx, err)
	} else {
		httpx.OkWithData(ctx, res)
	}
}

func (h *User) Info(ctx *gin.Context) {
	var req domain.IdPathReq
	if err := httpx.BindAndValidate(ctx, &req); err != nil {
		httpx.FailWithErr(ctx, err)
		return
	}

	res, err := h.user.Info(ctx.Request.Context(), &req)
	if err != nil {
		httpx.FailWithErr(ctx, err)
	} else {
		httpx.OkWithData(ctx, res)
	}
}

func (h *User) Create(ctx *gin.Context) {
	var req domain.User
	if err := httpx.BindAndValidate(ctx, &req); err != nil {
		httpx.FailWithErr(ctx, err)
		return
	}

	err := h.user.Create(ctx.Request.Context(), &req)
	if err != nil {
		httpx.FailWithErr(ctx, err)
	} else {
		httpx.Ok(ctx)
	}
}

func (h *User) Edit(ctx *gin.Context) {
	var req domain.User
	if err := httpx.BindAndValidate(ctx, &req); err != nil {
		httpx.FailWithErr(ctx, err)
		return
	}

	err := h.user.Edit(ctx.Request.Context(), &req)
	if err != nil {
		httpx.FailWithErr(ctx, err)
	} else {
		httpx.Ok(ctx)
	}
}

func (h *User) Delete(ctx *gin.Context) {
	var req domain.IdPathReq
	if err := httpx.BindAndValidate(ctx, &req); err != nil {
		httpx.FailWithErr(ctx, err)
		return
	}

	err := h.user.Delete(ctx.Request.Context(), &req)
	if err != nil {
		httpx.FailWithErr(ctx, err)
	} else {
		httpx.Ok(ctx)
	}
}

func (h *User) List(ctx *gin.Context) {
	var req domain.UserListReq
	if err := httpx.BindAndValidate(ctx, &req); err != nil {
		httpx.FailWithErr(ctx, err)
		return
	}

	res, err := h.user.List(ctx.Request.Context(), &req)
	if err != nil {
		httpx.FailWithErr(ctx, err)
	} else {
		httpx.OkWithData(ctx, res)
	}
}

func (h *User) UpPassword(ctx *gin.Context) {
	var req domain.UpPasswordReq
	if err := httpx.BindAndValidate(ctx, &req); err != nil {
		httpx.FailWithErr(ctx, err)
		return
	}

	err := h.user.UpPassword(ctx.Request.Context(), &req)
	if err != nil {
		httpx.FailWithErr(ctx, err)
	} else {
		httpx.Ok(ctx)
	}
}
