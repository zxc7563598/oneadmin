package role

import "github.com/zxc7563598/oneadmin/internal/service/role"

type Handler struct {
	roleSvc *role.Service
}

func New(roleSvc *role.Service) *Handler {
	return &Handler{
		roleSvc: roleSvc,
	}
}
