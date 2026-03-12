package response

import (
	"github.com/gin-gonic/gin"
	"github.com/zxc7563598/oneadmin/internal/i18n"
)

// Success 通用接口成功返回
func Success(c *gin.Context, data any) {
	lang, _ := c.Get("lang")
	c.JSON(200, gin.H{
		"code": 0,
		"msg":  i18n.GetMessage(lang.(string), 0),
		"data": data,
	})
}

// Error 通用接口失败返回
func Error(c *gin.Context, code int, data ...any) {
	lang, _ := c.Get("lang")
	var d any
	if len(data) > 0 {
		d = data[0]
	} else {
		d = nil
	}
	c.JSON(200, gin.H{
		"code": code,
		"msg":  i18n.GetMessage(lang.(string), code),
		"data": d,
	})
}
