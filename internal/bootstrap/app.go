package bootstrap

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/zxc7563598/oneadmin/internal/i18n"
	"gorm.io/gorm"
)

func NewApp(db *gorm.DB, rdb *redis.Client) *gin.Engine {
	// repository
	repos := InitRepositories(db)
	// service
	services := InitServices(repos, db, rdb)
	// handler
	handlers := InitHandlers(services)
	// i18n
	// 也可以通过 `i18n.Load("zh", "internal/i18n/locales/zh.yaml")` 加载单独的语言
	i18n.LoadLocales("internal/i18n/locales")
	// 注册路由
	r := gin.New()
	r = RouteRegister(r, rdb, handlers)
	return r
}
