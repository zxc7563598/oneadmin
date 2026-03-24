package bootstrap

import (
	"github.com/redis/go-redis/v9"
	"github.com/zxc7563598/oneadmin/internal/service/admin"
	"github.com/zxc7563598/oneadmin/internal/service/menu"
	"github.com/zxc7563598/oneadmin/internal/service/role"
	"gorm.io/gorm"
)

type Services struct {
	Admin admin.Service
	Role  role.Service
	Menu  menu.Service
}

func InitServices(repo *Repositories, db *gorm.DB, rdb *redis.Client) *Services {
	return &Services{
		Admin: *admin.New(repo.Admin, repo.AdminRole, repo.Role, db, rdb),
		Role:  *role.New(repo.Role, repo.Admin, repo.RoleMenu, repo.AdminRole, repo.Menu, db),
		Menu:  *menu.New(repo.Menu),
	}
}
