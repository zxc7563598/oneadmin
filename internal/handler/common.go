package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/zxc7563598/oneadmin/internal/validation"
	"go.uber.org/zap"
)

// GetAdminID 获取 JWT 携带的管理员 ID
func GetAdminID(c *gin.Context) (uint64, bool) {
	v, exists := c.Get("adminID")
	if !exists {
		return 0, false
	}
	adminID, ok := v.(uint64)
	return adminID, ok
}

// BindAndValidate 绑定请求参数并进行验证，失败将得到错误码
func BindAndValidate(c *gin.Context, req any) (int, bool, error) {
	if err := c.ShouldBindJSON(req); err != nil {
		code := validation.Parse(err, req)
		return code, false, err
	}
	return 0, true, nil
}

// ErrorLog 根据异常 Code 区分级别，封装异常日志信息
func ErrorLog(log *zap.Logger, msg string, code int, err error, fields ...zap.Field) {
	newFields := make([]zap.Field, 0, len(fields)+2)
	newFields = append(newFields, fields...)
	newFields = append(newFields, zap.Int("code", code))
	if err != nil {
		newFields = append(newFields, zap.Error(err))
	}
	switch {
	case code >= 50000:
		log.Error(msg, newFields...)
	case code >= 30000:
		log.Warn(msg, newFields...)
	default:
		log.Info(msg, newFields...)
	}
}
