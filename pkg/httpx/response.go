package httpx

import (
	"errors"
	"sync"

	"github.com/gin-gonic/gin"
)

const (
	ERROR   = 500
	SUCCESS = 200
)

var (
	ERRORMSG   = errors.New("fail")
	SUCCESSMSG = "success"
)

var (
	errorHandler func(ctx *gin.Context, err error) (int, error)
	errorLock    sync.RWMutex
	okHandler    func(ctx *gin.Context, data any) any
	okLock       sync.RWMutex
)

var NULL = map[string]interface{}{}

func SetOkHandler(handler func(ctx *gin.Context, data any) any) {
	okHandler = handler
}

func SetErrorHandler(handler func(ctx *gin.Context, err error) (int, error)) {
	errorHandler = handler
}

type Response struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
	Msg  string      `json:"msg"`
}

func Result(ctx *gin.Context, code int, data interface{}, msg string) {
	ctx.JSON(200, &Response{
		code,
		data,
		msg,
	})
}

func Ok(ctx *gin.Context) {
	OkWithData(ctx, NULL)
}

func OkWithData(ctx *gin.Context, data interface{}) {
	okLock.RLock()
	handler := okHandler
	okLock.RUnlock()
	if handler != nil {
		data = handler(ctx, data)
	}
	Result(ctx, SUCCESS, data, SUCCESSMSG)
}

func Fail(ctx *gin.Context) {
	FailWithErr(ctx, ERRORMSG)
}

func FailWithErr(ctx *gin.Context, err error) {
	errorLock.RLock()
	handler := errorHandler
	errorLock.RUnlock()

	code := ERROR
	if handler != nil {
		code, err = handler(ctx, err)
	}
	Result(ctx, code, NULL, err.Error())
}
