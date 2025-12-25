package utils

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
)

// JwtSecret 定义密钥（在生产环境中，可以从 config/env 中读取）
var JwtSecret = []byte("123456")

// Claims 自定义载荷，可以存一些非敏感信息，比如 UserID
type Claims struct {
	UserID uint `json:"user_id"`
	jwt.RegisteredClaims
}

func GenerateToken(userID uint) (string, error) {
	// 设置有效期24h
	expirationTime := time.Now().Add(24 * time.Hour)

	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime), // 过期时间
			IssuedAt:  jwt.NewNumericDate(time.Now()),     // 签发时间
			Issuer:    "go-link-project",                  // 签发人
		},
	}

	// 使用HS256算法创建一个token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// 生成签名字符串
	tokenString, err := token.SignedString(JwtSecret)
	return "Bearer " + tokenString, err
}

func ParseToken(tokenString string) (*Claims, error) {
	// jwt解析函数 参数表： 令牌字符串 用于接收解析后的JWT的payload声明的claims 一个回调函数，用于验证签名的秘钥 可变参数（例如时间验证、算法等等）
	// 返回解析后的token
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (any, error) {
		return JwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	// 没解析出来
	return nil, err
}
