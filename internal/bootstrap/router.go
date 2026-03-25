package bootstrap

import (
	"io"
	"io/fs"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/zxc7563598/oneadmin/internal/middleware"
	"github.com/zxc7563598/oneadmin/internal/webui"
)

func RouteRegister(r *gin.Engine, rdb *redis.Client, handlers *Handlers) *gin.Engine {
	r.RedirectTrailingSlash = false
	r.RedirectFixedPath = false
	// 日志注册
	if gin.Mode() != gin.ReleaseMode {
		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}
	// 中间件注册
	r.Use(gin.Logger(), gin.Recovery(), middleware.CORSMiddleware(), middleware.LocaleMiddleware())
	// web路由
	admin := r.Group("/admin")
	// web.Use(middleware.WebBasicAuth(webUser, webPass))
	RegisterWeb(admin)
	// api路由
	adminApi := r.Group("/api/admin")
	adminApi.POST("/auth/login", handlers.Admin.Login)
	adminApi.POST("/auth/switch-role", middleware.AdminAuth(rdb), handlers.Admin.SwitchRole)
	adminApi.POST("/auth/logout", middleware.AdminAuth(rdb), handlers.Admin.Logout)
	adminApi.POST("/auth/refresh", handlers.Admin.Refresh)
	adminApi.POST("/auth/detail", middleware.AdminAuth(rdb), handlers.Admin.Details)
	adminApi.POST("/admin/list", middleware.AdminAuth(rdb), handlers.Admin.ListPage)
	adminApi.POST("/admin/delete", middleware.AdminAuth(rdb), handlers.Admin.Delete)
	adminApi.POST("/admin/save", middleware.AdminAuth(rdb), handlers.Admin.Save)
	adminApi.POST("/admin/update-profile", middleware.AdminAuth(rdb), handlers.Admin.UpdateProfile)
	adminApi.POST("/auth/change-password", middleware.AdminAuth(rdb), handlers.Admin.ChangePassword)
	adminApi.POST("/admin/update-password", middleware.AdminAuth(rdb), handlers.Admin.UpdatePassword)
	adminApi.POST("/roles/list", middleware.AdminAuth(rdb), handlers.Role.ListPage)
	adminApi.POST("/roles/all", middleware.AdminAuth(rdb), handlers.Role.ListAll)
	adminApi.POST("/roles/save", middleware.AdminAuth(rdb), handlers.Role.Save)
	adminApi.POST("/roles/delete", middleware.AdminAuth(rdb), handlers.Role.Delete)
	adminApi.POST("/roles/add-role-users", middleware.AdminAuth(rdb), handlers.Role.AddRoleUsers)
	adminApi.POST("/roles/remove-role-users", middleware.AdminAuth(rdb), handlers.Role.RemoveRoleUsers)
	adminApi.POST("/roles/permissions", middleware.AdminAuth(rdb), handlers.Role.Permissions)
	adminApi.POST("/menu/list", middleware.AdminAuth(rdb), handlers.Menu.List)
	adminApi.POST("/menu/validate", middleware.AdminAuth(rdb), handlers.Menu.Validate)
	adminApi.POST("/menu/buttons", middleware.AdminAuth(rdb), handlers.Menu.Buttons)
	adminApi.POST("/menu/save", middleware.AdminAuth(rdb), handlers.Menu.Save)
	adminApi.POST("/menu/toggle", middleware.AdminAuth(rdb), handlers.Menu.Toggle)
	adminApi.POST("/menu/delete", middleware.AdminAuth(rdb), handlers.Menu.Delete)
	return r
}

func RegisterWeb(admin *gin.RouterGroup) {
	sub, err := fs.Sub(webui.Dist, "dist")
	if err != nil {
		panic(err)
	}
	fileServer := http.FileServer(http.FS(sub))
	admin.GET("", func(c *gin.Context) {
		c.Redirect(http.StatusFound, "/admin/")
	})
	admin.GET("/*filepath", func(c *gin.Context) {
		path := c.Param("filepath")
		if len(path) > 0 && path[0] == '/' {
			path = path[1:]
		}
		if _, err := sub.Open(path); err == nil {
			http.StripPrefix("/admin/", fileServer).ServeHTTP(c.Writer, c.Request)
			return
		}
		index, err := sub.Open("index.html")
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return
		}
		defer index.Close()
		c.Header("Content-Type", "text/html; charset=utf-8")
		c.Status(http.StatusOK)
		io.Copy(c.Writer, index)
	})
}
