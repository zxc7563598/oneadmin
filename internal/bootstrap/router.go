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
	r.Use(gin.Logger(), gin.Recovery(), middleware.CORSMiddleware(), middleware.LocaleMiddleware())
	// web路由
	admin := r.Group("/admin")
	// web.Use(middleware.WebBasicAuth(webUser, webPass))
	RegisterWeb(admin)
	// api路由
	adminApi := r.Group("/api/admin")
	adminApi.POST("/auth/login", handlers.Admin.Login)                                                  // 完成
	adminApi.POST("/auth/switch-role", middleware.AdminAuth(rdb), handlers.Admin.SwitchRole)            // 完成
	adminApi.POST("/auth/logout", middleware.AdminAuth(rdb), handlers.Admin.Logout)                     // 完成
	adminApi.POST("/auth/refresh", handlers.Admin.Refresh)                                              // 完成
	adminApi.POST("/auth/detail", middleware.AdminAuth(rdb), handlers.Admin.Details)                    // 完成
	adminApi.POST("/admin/list", middleware.AdminAuth(rdb), handlers.Admin.ListPage)                    // 完成
	adminApi.POST("/admin/delete", middleware.AdminAuth(rdb), handlers.Admin.Delete)                    // 完成
	adminApi.POST("/admin/save", middleware.AdminAuth(rdb), handlers.Admin.Save)                        // 完成
	adminApi.POST("/admin/update-profile", middleware.AdminAuth(rdb), handlers.Admin.UpdateProfile)     // 完成
	adminApi.POST("/auth/change-password", middleware.AdminAuth(rdb), handlers.Admin.ChangePassword)    // 完成
	adminApi.POST("/admin/update-password", middleware.AdminAuth(rdb), handlers.Admin.UpdatePassword)   // 完成
	adminApi.POST("/roles/list", middleware.AdminAuth(rdb), handlers.Role.ListPage)                     // 完成
	adminApi.POST("/roles/all", middleware.AdminAuth(rdb), handlers.Role.ListAll)                       // 完成
	adminApi.POST("/roles/save", middleware.AdminAuth(rdb), handlers.Role.Save)                         // 完成
	adminApi.POST("/roles/delete", middleware.AdminAuth(rdb), handlers.Role.Delete)                     // 完成
	adminApi.POST("/roles/add-role-users", middleware.AdminAuth(rdb), handlers.Role.AddRoleUsers)       // 完成
	adminApi.POST("/roles/remove-role-users", middleware.AdminAuth(rdb), handlers.Role.RemoveRoleUsers) // 完成
	adminApi.POST("/roles/permissions", middleware.AdminAuth(rdb), handlers.Role.Permissions)           // 完成

	adminApi.POST("/menu/list", middleware.AdminAuth(rdb), handlers.Menu.List)         // 完成
	adminApi.POST("/menu/validate", middleware.AdminAuth(rdb), handlers.Menu.Validate) // 完成
	adminApi.POST("/menu/buttons", middleware.AdminAuth(rdb), handlers.Menu.Buttons)   // 完成
	adminApi.POST("/menu/save", middleware.AdminAuth(rdb), handlers.Menu.Save)         // 完成
	adminApi.POST("/menu/toggle", middleware.AdminAuth(rdb), handlers.Menu.Toggle)     // 完成
	adminApi.POST("/menu/delete", middleware.AdminAuth(rdb), handlers.Menu.Delete)     // 完成
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
