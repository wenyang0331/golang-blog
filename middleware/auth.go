package middleware

import (
	"homework_blog/utils"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func Auth(jwtSecret []byte) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从 Header 获取 Token
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.Error(c, http.StatusUnauthorized, "Authorization header required")
			c.Abort()
			return
		}

		// 提取 Token（Bearer <token>）
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			utils.Error(c, http.StatusUnauthorized, "Invalid authorization header format")
			c.Abort()
			return
		}

		tokenString := parts[1]

		// 诊断：打印 token 基本信息
		log.Printf("[JWT] Received token, length=%d, prefix=%.20q", len(tokenString), tokenString)

		// 验证 Token
		claims, err := utils.ParseToken(tokenString, jwtSecret)
		if err != nil {
			log.Printf("[JWT] Token validation failed: %v", err)
			utils.Error(c, http.StatusUnauthorized, "Invalid token: "+err.Error())
			c.Abort()
			return
		}

		// 将用户信息存储到 Context
		c.Set("userID", claims.UserID)
		c.Set("username", claims.Username)

		c.Next()
	}
}
