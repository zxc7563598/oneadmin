package role_menu

import (
	"github.com/zxc7563598/oneadmin/internal/model"
	"github.com/zxc7563598/oneadmin/internal/repository/base"
)

type Repository interface {
	base.Repository[model.RoleMenu]
}
