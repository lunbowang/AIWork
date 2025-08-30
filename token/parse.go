package token

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt"
)

// JWT标准声明字段常量定义，对应JWT规范中的标准Claims字段
const (
	jwtAudience   = "aud"           // 受众（接收者）
	jwtExpire     = "exp"           // 过期时间戳
	jwtId         = "jti"           // JWT唯一标识
	jwtIssueAt    = "iat"           // 签发时间戳
	jwtIssuer     = "iss"           // 签发者
	jwtNotBefore  = "nbf"           // 生效时间（在此时间前令牌无效）
	jwtSubject    = "sub"           // 主题
	Authorization = "Authorization" // HTTP请求头中存储令牌的键名
)

// 令牌处理相关的错误变量定义
var (
	ErrTokenNotFound = errors.New("不存在token")
	ErrTokenInvalid  = errors.New("token is invalid")
	ErrClaimsInvalid = errors.New("invalid token claims")
)

// Parse 用于解析和验证JWT令牌的结构体
// 包含解析令牌所需的密钥信息
type Parse struct {
	AccessSecret string
}

// NewTokenParse 创建一个新的Parse实例
func NewTokenParse(secret string) *Parse {
	return &Parse{AccessSecret: secret}
}

// Parse 从HTTP请求中提取并解析JWT令牌
func (p *Parse) Parse(r *http.Request) (jwt.MapClaims, string, error) {
	tokenStr := p.extractTokenFromHeader(r)
	if len(tokenStr) == 0 {
		return nil, tokenStr, ErrTokenNotFound
	}
	return p.ParseToken(tokenStr)
}

// ParseToken 解析并验证JWT令牌字符串
func (p *Parse) ParseToken(tokenStr string) (jwt.MapClaims, string, error) {
	// 解析令牌，使用提供的密钥验证签名
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte(p.AccessSecret), nil
	})
	if err != nil {
		return nil, tokenStr, err
	}
	// 校验token是否有效
	if !token.Valid {
		return nil, tokenStr, ErrTokenInvalid
	}

	// 将令牌声明转换为MapClaims类型
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, tokenStr, ErrClaimsInvalid
	}

	// 验证声明的有效性（如过期时间等）
	if err = claims.Valid(); err != nil {
		return nil, tokenStr, err
	}
	return claims, tokenStr, nil
}

// ParseWithContext 从HTTP请求中解析令牌，并将令牌信息存入请求上下文
func (p *Parse) ParseWithContext(r *http.Request) (*http.Request, error) {
	// 解析请求中的令牌
	claims, tokenStr, err := p.Parse(r)
	if err != nil {
		return r, err
	}

	// 获取请求的上下文
	ctx := r.Context()

	// 将自定义声明字段存入上下文（忽略标准字段）
	for k, v := range claims {
		switch k {
		case jwtAudience, jwtExpire, jwtId, jwtIssueAt, jwtIssuer, jwtNotBefore, jwtSubject:
			// 忽略标准声明字段
		default:
			ctx = context.WithValue(ctx, k, v)
		}
	}

	// 将令牌字符串存入上下文，键为Authorization
	ctx = context.WithValue(ctx, Authorization, tokenStr)

	// 返回带有新上下文的请求对象
	return r.WithContext(ctx), nil
}

// extractTokenFromHeader 从HTTP请求头中提取JWT令牌字符串
func (p *Parse) extractTokenFromHeader(r *http.Request) string {
	// 从请求头中获取Authorization字段值
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return ""
	}

	// 按空格分割头信息，期望格式为"Bearer <token>"
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return authHeader
	}

	// 返回分割后的令牌部分
	return parts[1]
}

// GetTokenStr 从上下文中获取存储的令牌字符串
func GetTokenStr(ctx context.Context) string {
	tokenStr, ok := ctx.Value(Authorization).(string)
	if !ok {
		return ""
	}
	return tokenStr
}

// VerifyJWTToken 从HTTP请求中提取并验证JWT令牌
// 参数：
//
//	secretKey: 用于验证令牌签名的密钥
//	r: HTTP请求对象，需包含Authorization头
//
// 返回值：
//
//	jwt.MapClaims: 解析后的令牌声明数据
//	error: 验证过程中出现的错误（如令牌不存在、无效等）
func VerifyJWTToken(secretKey string, r *http.Request) (jwt.MapClaims, error) {
	// 创建Parse实例，用于处理令牌解析
	parser := NewTokenParse(secretKey)

	// 从请求中解析令牌并验证
	claims, _, err := parser.Parse(r)
	if err != nil {
		// 根据不同错误类型返回更具体的错误信息
		switch {
		case errors.Is(err, ErrTokenNotFound):
			return nil, fmt.Errorf("验证失败: %w", ErrTokenNotFound)
		case errors.Is(err, ErrTokenInvalid):
			return nil, fmt.Errorf("验证失败: %w", ErrTokenInvalid)
		case errors.Is(err, ErrClaimsInvalid):
			return nil, fmt.Errorf("验证失败: %w", ErrClaimsInvalid)
		default:
			return nil, fmt.Errorf("验证失败: %w", err)
		}
	}

	// 验证成功，返回解析后的声明数据
	return claims, nil
}
