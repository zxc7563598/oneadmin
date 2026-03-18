package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zxc7563598/oneadmin/internal/i18n"
)

func respond(c *gin.Context, code int, lang string, data any) {
	if lang == "" {
		lang = i18n.GetLang(c.Request.Context())
	}
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  i18n.GetMessage(lang, code),
		"data": data,
	})
}

// Success 接口成功返回
func Success(c *gin.Context, lang string, data any) {
	respond(c, 0, lang, data)
}

// Error 接口失败返回
func Error(c *gin.Context, lang string, code int, data ...any) {
	d := any(nil)
	if len(data) > 0 {
		d = data[0]
	}
	respond(c, code, lang, d)
}
