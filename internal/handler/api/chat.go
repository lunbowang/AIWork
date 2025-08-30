package api

import (
	"github.com/gin-gonic/gin"

	"ai/internal/domain"
	"ai/internal/logic"
	"ai/internal/svc"
	"ai/pkg/httpx"
)

type Chat struct {
	svcCtx *svc.ServiceContext
	chat   logic.Chat
}

func NewChat(svcCtx *svc.ServiceContext, chat logic.Chat) *Chat {
	return &Chat{
		svcCtx: svcCtx,
		chat:   chat,
	}
}

func (h *Chat) InitRegister(engine *gin.Engine) {
	g := engine.Group("v1/chat", h.svcCtx.Jwt.Handler)
	g.POST("", h.Chat)
}

func (h *Chat) Chat(ctx *gin.Context) {
	var req domain.ChatReq
	if err := httpx.BindAndValidate(ctx, &req); err != nil {
		httpx.FailWithErr(ctx, err)
		return
	}

	res, err := h.chat.AIChat(ctx.Request.Context(), &req)
	if err != nil {
		httpx.FailWithErr(ctx, err)
	} else {
		httpx.OkWithData(ctx, res)
	}
}
