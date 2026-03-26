package bootstrap

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/zxc7563598/oneadmin/internal/config"
	"github.com/zxc7563598/oneadmin/internal/i18n"
	"github.com/zxc7563598/oneadmin/internal/logger"
	"github.com/zxc7563598/oneadmin/internal/migrate"
	"github.com/zxc7563598/oneadmin/pkg/jwt"
)

func NewApp(cfg *config.Config) *gin.Engine {
	// 初始化日志
	logger.InitAll()
	// 初始化redis
	rdb, err := config.InitRedis(cfg)
	if err != nil {
		log.Fatalf("Redis配置存在但连接失败: %v", err)
	}
	if rdb != nil {
		defer rdb.Close()
		log.Println("Redis已启用")
	}
	// 初始化数据库
	db, err := config.InitDB(cfg)
	if err != nil {
		log.Fatalf("无法初始化数据库: %v", err)
	}
	// 自动建表
	if err := migrate.Run(db); err != nil {
		panic(err)
	}
	// 初始化填充数据
	if err := migrate.Seed(db); err != nil {
		panic(err)
	}
	// 初始化jwt
	jwt.Init(cfg.JWT)
	// 处理依赖注入
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
