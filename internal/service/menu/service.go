package menu

import "github.com/zxc7563598/oneadmin/internal/repository/menu"

type Service struct {
	menuRepo menu.Repository
}

func New(menuRepo menu.Repository) *Service {
	return &Service{
		menuRepo: menuRepo,
	}
}
