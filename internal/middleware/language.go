package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
)

// LocaleMiddleware 多语言中间件
func LocaleMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		lang := c.Query("lang")
		if lang == "" {
			lang = c.GetHeader("X-Lang")
		}
		if lang == "" {
			lang = c.GetHeader("Accept-Language")
		}
		if lang == "" {
			lang = "zh"
		}
		// zh-CN -> zh
		if strings.Contains(lang, "-") {
			lang = strings.Split(lang, "-")[0]
		}
		c.Set("lang", strings.ToLower(lang))
		c.Next()
	}
}
