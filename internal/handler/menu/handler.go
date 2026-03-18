package menu

import "github.com/zxc7563598/oneadmin/internal/service/menu"

type Handler struct {
	menuSvc *menu.Service
}

func New(menuSvc *menu.Service) *Handler {
	return &Handler{
		menuSvc: menuSvc,
	}
}
