package role

import "github.com/zxc7563598/oneadmin/internal/repository/role"

type Service struct {
	roleRepo role.Repository
}

func New(roleRepo role.Repository) *Service {
	return &Service{
		roleRepo: roleRepo,
	}
}
