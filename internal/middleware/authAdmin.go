package middleware

import (
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"

	"github.com/zxc7563598/oneadmin/internal/response"
	"github.com/zxc7563598/oneadmin/pkg/jwt"
)

const RenewThreshold = 15 * time.Minute

func AdminAuth(rdb *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Error(c, 20001)
			return
		}
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.Error(c, 20002)
			return
		}
		tokenString := parts[1]
		claims, err := jwt.ParseToken(tokenString)
		if err != nil {
			response.Error(c, 20003)
			return
		}
		if claims.Type != "access" {
			response.Error(c, 20004)
			return
		}
		// Redis 单点登录校验
		if rdb != nil {
			ctx := c.Request.Context()
			key := fmt.Sprintf("admin:token:%d", claims.ID)
			redisToken, err := rdb.Get(ctx, key).Result()
			if err == redis.Nil {
				response.Error(c, 20005)
				return
			}
			if err != nil {
				response.Error(c, 50001)
				return
			}
			if redisToken != tokenString {
				response.Error(c, 20006)
				return
			}
		}
		expireTime := claims.ExpiresAt.Time
		remain := time.Until(expireTime)
		// 滑动续期
		if remain < RenewThreshold {
			newToken, err := jwt.GenerateAccessToken(
				claims.ID,
				"admin",
			)
			if err == nil {
				c.Header("X-New-Token", newToken)
				if rdb != nil {
					ctx := c.Request.Context()
					key := fmt.Sprintf("admin:token:%d", claims.ID)
					rdb.Set(
						ctx,
						key,
						newToken,
						jwt.AccessTTL(),
					)
				}
			}
		}
		// 写入上下文
		c.Set("admin_id", claims.ID)
		c.Next()
	}
}
