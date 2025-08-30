package handler

import (
	"github.com/gin-gonic/gin"
)

type cause interface {
	Cause() error
}

func ErrorHandler(ctx *gin.Context, err error) (int, error) {
	var e error
	if ce, ok := err.(cause); ok {
		e = ce.Cause()
	} else {
		e = err
	}

	return 500, e
}
