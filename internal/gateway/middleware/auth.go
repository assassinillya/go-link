package middleware

import (
	"github.com/gin-gonic/gin"
	"go-link/pkg/utils"
	"net/http"
	"strings"
)

// JWTAuth jwt中间件：检查请求头里有没有合法的Token
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取请求头 Authorization
		// 格式通常是 Bearer xxxxx.yyyyy.zzzzz
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "请求未携带Token"})
			c.Abort() // 阻止请求继续往下走
			return
		}

		// 截取真正的Token部分
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token 格式错误 (Bearer type)"})
			c.Abort()
			return
		}

		// parts[1]为token部分
		// 开始解析 Token
		claims, err := utils.ParseToken(parts[1])
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token 无效或已过期"})
			c.Abort()
			return
		}

		// 将userID存入ctx中
		// 接下来后续处理函数就能凭token知道是什么用户在请求了
		c.Set("userID", claims.UserID)

		// 放行
		c.Next()
	}
}
