package token

import (
	"context"

	"github.com/golang-jwt/jwt"
)

const Identify = "LunBoWang"

// GetJwtToken 生成Token
func GetJwtToken(secretKey string, iat, seconds int64, uid string) (string, error) {
	// 创建JWT声明（claims），用于存储自定义数据和标准字段
	claims := make(jwt.MapClaims)
	// 设置令牌过期时间：签发时间+有效期
	claims["exp"] = iat + seconds
	// 设置令牌签发时间
	claims["iat"] = iat
	// 存储用户唯一标识，键名为全局常量Identify
	claims[Identify] = uid
	// 创建一个使用HS256算法的JWT令牌实例
	token := jwt.New(jwt.SigningMethodHS256)
	// 将声明设置到令牌中
	token.Claims = claims
	// 使用密钥对令牌进行签名，生成最终的令牌字符串
	return token.SignedString([]byte(secretKey))
}

// GetUId 从context中获取用户ID
func GetUId(ctx context.Context) string {
	var uid string
	// 从context中获取键为Identify的值，并尝试转换为string类型
	if jsonUid, ok := ctx.Value(Identify).(string); ok {
		uid = jsonUid
	}
	return uid
}
