package middleware

import (
	"ai/pkg/httpx"
	"ai/token"

	"github.com/gin-gonic/gin"
)

// Jwt 封装了JWT令牌解析器的结构体，用于处理HTTP请求中的JWT验证
// 作为Gin框架的中间件使用，负责从请求中提取并验证JWT令牌
type Jwt struct {
	// tokenParser 用于解析和验证JWT令牌的解析器实例
	tokenParser *token.Parse
}

// NewJwt 创建一个新的Jwt实例
// 参数secret: 用于验证JWT签名的密钥
// 返回值: 初始化后的Jwt指针，包含令牌解析器
func NewJwt(secret string) *Jwt {
	return &Jwt{
		// 初始化令牌解析器，传入签名密钥
		tokenParser: token.NewTokenParse(secret),
	}
}

// Handler 实现Gin框架的中间件接口，用于处理请求中的JWT令牌验证
// 将令牌解析后的数据存入请求上下文，供后续处理函数使用
func (m *Jwt) Handler(ctx *gin.Context) {
	// 调用令牌解析器的ParseWithContext方法，从请求中解析令牌并将信息存入上下文
	// 该方法会处理Authorization请求头的提取、令牌验证和上下文注入
	r, err := m.tokenParser.ParseWithContext(ctx.Request)
	if err != nil {
		httpx.FailWithErr(ctx, err)
		// 终止后续请求处理流程，不再调用后续中间件和处理函数
		ctx.Abort()
		return
	}

	// 解析成功，更新请求对象为包含新上下文的请求
	ctx.Request = r
	// 调用Next()方法，将请求传递给下一个中间件或处理函数
	ctx.Next()
}
