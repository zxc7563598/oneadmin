package bootstrap

import (
	"github.com/zxc7563598/oneadmin/internal/handler/admin"
	"github.com/zxc7563598/oneadmin/internal/handler/menu"
	"github.com/zxc7563598/oneadmin/internal/handler/role"
)

type Handlers struct {
	Admin *admin.Handler
	Menu  *menu.Handler
	Role  *role.Handler
}

func InitHandlers(svc *Services) *Handlers {
	return &Handlers{
		Admin: admin.New(&svc.Admin),
		Menu:  menu.New(&svc.Menu),
		Role:  role.New(&svc.Role),
	}
}
