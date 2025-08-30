package middleware

import (
	"fmt"
	"time"

	"gitee.com/dn-jinmin/tlog"
	"github.com/gin-gonic/gin"
)

type Log struct{}

func NewLog() *Log {
	return &Log{}
}

func (w *Log) Handler(ctx *gin.Context) {
	startTime := time.Now()
	url := fmt.Sprintf("%s:%s", ctx.Request.URL.Path, ctx.Request.Method)

	ctx.Request = ctx.Request.WithContext(tlog.TraceStart(ctx.Request.Context()))
	defer func() {
		tlog.InfoCtx(ctx.Request.Context(), url, "time", tlog.RTField(startTime, time.Now()))
	}()

	ctx.Next()
}
