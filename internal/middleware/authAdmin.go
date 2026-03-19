package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"

	"github.com/zxc7563598/oneadmin/internal/response"
	"github.com/zxc7563598/oneadmin/pkg/jwt"
)

func AdminAuth(rdb *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Error(c, "", 10104)
			return
		}
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.Error(c, "", 10105)
			return
		}
		tokenString := parts[1]
		claims, err := jwt.ParseToken(tokenString)
		if err != nil {
			response.Error(c, "", 10101)
			return
		}
		if claims.Type != "access" {
			response.Error(c, "", 10102)
			return
		}
		// Redis 单点登录校验
		// 如果未启动 Redis 则不进行单点校正，用户会在 AccessToken 过期后，重新申请 RefreshToken 时才会被踢出
		if rdb != nil {
			ctx := c.Request.Context()
			key := jwt.AdminTokenKey(claims.ID)
			redisToken, err := rdb.Get(ctx, key).Result()
			if err == redis.Nil {
				response.Error(c, "", 10106)
				return
			}
			if err != nil {
				response.Error(c, "", 10107)
				return
			}
			if redisToken != tokenString {
				response.Error(c, "", 20001)
				return
			}
		}
		// 写入上下文
		c.Set("adminID", claims.ID)
		c.Set("roleID", claims.RoleID)
		c.Next()
	}
}
