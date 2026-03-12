package bootstrap

import (
	"github.com/gin-gonic/gin"
	"github.com/zxc7563598/oneadmin/internal/i18n"
	"gorm.io/gorm"
)

func NewApp(db *gorm.DB) *gin.Engine {
	// i18n
	// 也可以通过 `i18n.Load("zh", "internal/i18n/locales/zh.yaml")` 加载单独的语言
	i18n.LoadLocales("internal/i18n/locales")
	// 注册路由
	r := gin.New()
	return r
}
