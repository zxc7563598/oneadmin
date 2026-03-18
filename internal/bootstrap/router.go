package bootstrap

import (
	"io/fs"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/zxc7563598/oneadmin/internal/middleware"
	"github.com/zxc7563598/oneadmin/internal/webui"
)

func RouteRegister(r *gin.Engine, rdb *redis.Client, handlers *Handlers) *gin.Engine {
	r.RedirectTrailingSlash = false
	r.RedirectFixedPath = false
	// 中间件注册
	r.Use(gin.Logger(), gin.Recovery(), middleware.LocaleMiddleware())
	// web路由
	admin := r.Group("/admin")
	// web.Use(middleware.WebBasicAuth(webUser, webPass))
	RegisterWeb(admin)
	// api路由
	adminApi := r.Group("/api/admin")
	adminApi.POST("/auth/login", handlers.Admin.Login)
	adminApi.POST("/auth/refresh", handlers.Admin.Refresh)
	adminApi.POST("/auth/logout", middleware.AdminAuth(rdb), handlers.Admin.Logout)
	adminApi.POST("/auth/switch-role", middleware.AdminAuth(rdb), handlers.Admin.SwitchRole)
	adminApi.POST("/auth/change-password", middleware.AdminAuth(rdb), handlers.Admin.ChangePassword)
	adminApi.POST("/admin/list", middleware.AdminAuth(rdb), handlers.Admin.ListPage)
	adminApi.POST("/admin/detail", middleware.AdminAuth(rdb), handlers.Admin.Details)
	adminApi.POST("/admin/save", middleware.AdminAuth(rdb), handlers.Admin.Save)
	adminApi.POST("/admin/delete", middleware.AdminAuth(rdb), handlers.Admin.Delete)
	adminApi.POST("/admin/update-password", middleware.AdminAuth(rdb), handlers.Admin.UpdatePassword)
	adminApi.POST("/admin/update-profile", middleware.AdminAuth(rdb), handlers.Admin.UpdateProfile)
	return r
}

func RegisterWeb(admin *gin.RouterGroup) {
	sub, err := fs.Sub(webui.Dist, "dist")
	if err != nil {
		panic(err)
	}
	fileServer := http.StripPrefix(
		"/admin/",
		http.FileServer(http.FS(sub)),
	)
	// /admin → /admin/
	admin.GET("", func(c *gin.Context) {
		c.Redirect(http.StatusFound, "/admin/")
	})
	// /admin/* 交给 FileServer
	admin.GET("/*filepath", func(c *gin.Context) {
		fileServer.ServeHTTP(c.Writer, c.Request)
	})
}
