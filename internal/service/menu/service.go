package menu

import (
	"context"
	"sort"

	"github.com/zxc7563598/oneadmin/internal/enum"
	"github.com/zxc7563598/oneadmin/internal/model"
	"github.com/zxc7563598/oneadmin/internal/repository/menu"
)

type Service struct {
	menuRepo menu.Repository
}

func New(menuRepo menu.Repository) *Service {
	return &Service{
		menuRepo: menuRepo,
	}
}

// MenuTree 用于获取全部菜单
func (s *Service) MenuTree(ctx context.Context) ([]MenuItem, int, error) {
	// 获取菜单信息
	menus, err := s.menuRepo.FindAll(ctx, nil)
	if err != nil {
		return nil, 60301, err
	}
	// 整理菜单信息
	list := make([]MenuItem, 0, len(menus))
	for _, v := range menus {
		list = append(list, MenuItem{
			ID:        v.ID,
			Code:      v.Code,
			Enable:    v.Enable == enum.EnableEnable,
			Show:      v.Show == enum.Yes,
			KeepAlive: v.KeepAlive == enum.Yes,
			Layout:    v.Layout,
			Type:      string(v.Type),
			ParentID:  v.ParentID,
			Name:      v.Name,
			Icon:      v.Icon,
			Path:      v.Path,
			Component: v.Component,
			Order:     v.Order,
		})
	}
	// 返回权限树
	return s.buildTree(list, 0), 0, nil
}

// MenuExists 用于获取菜单是否存在
func (s *Service) MenuExists(ctx context.Context, path string) (bool, int, error) {
	has, err := s.menuRepo.PathToExists(ctx, nil, path)
	if err != nil {
		return false, 60301, err
	}
	return has, 0, nil
}

// MenuButtons 用于获取菜单下的按钮
func (s *Service) MenuButtons(ctx context.Context, parentID uint64) ([]MenuItem, int, error) {
	buttons, err := s.menuRepo.GetMenuButtons(ctx, nil, parentID)
	if err != nil {
		return nil, 60301, err
	}
	// 整理菜单信息
	list := make([]MenuItem, 0, len(buttons))
	for _, v := range buttons {
		list = append(list, MenuItem{
			ID:        v.ID,
			Code:      v.Code,
			Enable:    v.Enable == enum.EnableEnable,
			Show:      v.Show == enum.Yes,
			KeepAlive: v.KeepAlive == enum.Yes,
			Layout:    v.Layout,
			Type:      string(v.Type),
			ParentID:  v.ParentID,
			Name:      v.Name,
			Icon:      v.Icon,
			Path:      v.Path,
			Component: v.Component,
			Order:     v.Order,
		})
	}
	return list, 0, nil
}

// Save 用于添加或变更数据
func (s *Service) Save(ctx context.Context, req SaveReq) (int, error) {
	enable := enum.EnableDisable
	if req.Enable {
		enable = enum.EnableEnable
	}
	show := enum.No
	if req.Enable {
		show = enum.Yes
	}
	keepAlive := enum.No
	if req.Enable {
		keepAlive = enum.Yes
	}
	if req.ID == nil {
		// 添加数据
		_, err := s.menuRepo.Create(ctx, nil, &model.Menu{
			Code:      req.Code,
			Enable:    enable,
			Show:      show,
			KeepAlive: keepAlive,
			Layout:    req.Layout,
			Type:      enum.MenuType(req.Type),
			ParentID:  req.ParentID,
			Name:      req.Name,
			Icon:      req.Icon,
			Path:      req.Path,
			Component: req.Component,
			Order:     req.Order,
		})
		if err != nil {
			return 60302, err
		}
	} else {
		// 变更数据
		if err := s.menuRepo.UpdateInfo(ctx, nil, *req.ID, model.MenuUpdateInfoForm{
			Code:      req.Code,
			Enable:    enable,
			Show:      show,
			KeepAlive: keepAlive,
			Layout:    req.Layout,
			Type:      enum.MenuType(req.Type),
			ParentID:  req.ParentID,
			Name:      req.Name,
			Icon:      req.Icon,
			Path:      req.Path,
			Component: req.Component,
			Order:     req.Order,
		}); err != nil {
			return 60303, err
		}
	}
	// 返回成功
	return 0, nil
}

// SetMenuEnable 用于切换菜单启动状态
func (s *Service) SetMenuEnable(ctx context.Context, id uint64) (int, error) {
	err := s.menuRepo.UpdateEnableToggle(ctx, nil, id)
	if err != nil {
		return 60304, err
	}
	return 0, nil
}

// Delete 用于删除菜单
func (s *Service) Delete(ctx context.Context, id uint64) (int, error) {
	err := s.menuRepo.Delete(ctx, nil, id)
	if err != nil {
		return 60305, err
	}
	return 0, nil
}

// buildTree 用于获取菜单权限树
func (s *Service) buildTree(list []MenuItem, parentID uint64) []MenuItem {
	// 构建 parent -> children map
	childrenMap := make(map[uint64][]MenuItem)
	for _, el := range list {
		childrenMap[el.ParentID] = append(childrenMap[el.ParentID], el)
	}
	// 递归
	var build func(pid uint64) []MenuItem
	build = func(pid uint64) []MenuItem {
		branch := childrenMap[pid]
		// 排序
		sort.Slice(branch, func(i, j int) bool {
			return branch[i].Order < branch[j].Order
		})
		// 构建 children
		for i := range branch {
			branch[i].Children = build(branch[i].ID)
		}
		return branch
	}
	return build(parentID)
}
