package httpx

import "github.com/gin-gonic/gin"

func BindAndValidate(ctx *gin.Context, v any) error {
	if err := ctx.ShouldBind(v); err != nil {
		return err
	}

	if err := ctx.ShouldBindUri(v); err != nil {
		return err
	}

	return nil
}
