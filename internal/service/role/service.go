package role

import (
	"context"
	"errors"
	"sort"

	"github.com/zxc7563598/oneadmin/internal/enum"
	"github.com/zxc7563598/oneadmin/internal/model"
	"github.com/zxc7563598/oneadmin/internal/repository/menu"
	"github.com/zxc7563598/oneadmin/internal/repository/role"
	"github.com/zxc7563598/oneadmin/internal/repository/role_menu"
	"gorm.io/gorm"
)

type Service struct {
	roleRepo     role.Repository
	roleMenuRepo role_menu.Repository
	menuRepo     menu.Repository
	db           *gorm.DB
}

const RoleCodeSuperAdmin = "SUPER_ADMIN"

func New(roleRepo role.Repository, roleMenuRepo role_menu.Repository, menuRepo menu.Repository, db *gorm.DB) *Service {
	return &Service{
		roleRepo:     roleRepo,
		roleMenuRepo: roleMenuRepo,
		menuRepo:     menuRepo,
		db:           db,
	}
}

// ListPage 用于获取角色分页信息
func (s *Service) ListPage(ctx context.Context, req ListPageReq) (ListPageResp, int, error) {
	// 获取列表数据
	offset, limit := req.OffsetLimit()
	roles, total, err := s.roleRepo.ListPage(ctx, nil, model.RoleListQuery{
		Name:   req.Name,
		Enable: req.Enable,
		Offset: offset,
		Limit:  limit,
	})
	if err != nil {
		return ListPageResp{}, 60201, err
	}
	// 获取菜单
	roleIDs := make([]uint64, 0, len(roles))
	for _, v := range roles {
		roleIDs = append(roleIDs, v.ID)
	}
	menus, err := s.roleMenuRepo.GetByRoleIDs(ctx, nil, roleIDs)
	menuIDs := make(map[uint64][]uint64)
	for _, v := range menus {
		menuIDs[v.RoleID] = append(menuIDs[v.RoleID], v.MenuID)
	}
	// 组装数据
	list := make([]ListPageItem, 0, len(roles))
	for _, v := range roles {
		list = append(list, ListPageItem{
			ID:            v.ID,
			Code:          v.Code,
			Name:          v.Name,
			Enable:        v.Enable == enum.EnableEnable,
			PermissionIds: menuIDs[v.ID],
		})
	}
	// 返回数据
	return ListPageResp{
		Total:    total,
		PageData: list,
	}, 0, nil
}

// ListAll 用于获取角色全部信息
func (s *Service) ListAll(ctx context.Context) ([]ListAllResp, int, error) {
	// 获取角色
	roles, err := s.roleRepo.FindAll(ctx, nil)
	if err != nil {
		return nil, 60201, err
	}
	// 组装数据
	list := make([]ListAllResp, 0, len(roles))
	for _, v := range roles {
		list = append(list, ListAllResp{
			ID:     v.ID,
			Code:   v.Code,
			Name:   v.Name,
			Enable: v.Enable == enum.EnableEnable,
		})
	}
	return list, 0, nil
}

// Save 用于创建或修改角色信息
func (s *Service) Save(ctx context.Context, req SaveReq) (int, error) {
	// 开启事务
	err := s.db.Transaction(func(tx *gorm.DB) error {
		// 变更角色信息
		roleID, err := s.saveRole(ctx, tx, req)
		if err != nil {
			return err
		}
		// 重新绑定角色对应的权限
		if req.MenuIDs != nil {
			if err := s.resetRoleMenus(ctx, tx, roleID, *req.MenuIDs); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return 60202, err
	}
	return 0, nil
}

// saveRole 用于变更角色信息
func (s *Service) saveRole(ctx context.Context, tx *gorm.DB, req SaveReq) (uint64, error) {

	enable := enum.EnableDisable
	if req.Enable != nil {
		if *req.Enable {
			enable = enum.EnableEnable
		}
	}
	if req.ID == nil {
		// 新增
		if req.Code == nil || req.Name == nil {
			return 0, errors.New("missing required fields")
		}
		role, err := s.roleRepo.Create(ctx, tx, &model.Role{
			Code:   *req.Code,
			Name:   *req.Name,
			Enable: enable,
		})
		if err != nil {
			return 0, err
		}
		return role.ID, nil
	}
	// 更新
	roleID := *req.ID
	if err := s.roleRepo.UpdateByID(ctx, tx, roleID, model.RoleForm{
		Code:   req.Code,
		Name:   req.Name,
		Enable: enable,
	}); err != nil {
		return 0, err
	}
	return roleID, nil
}

// resetRoleMenus 用于重新绑定角色对应的权限
func (s *Service) resetRoleMenus(ctx context.Context, tx *gorm.DB, roleID uint64, menuIDs []uint64) error {
	// 删除旧关系
	if err := s.roleMenuRepo.DeleteByRoleID(ctx, tx, roleID); err != nil {
		return err
	}
	// 没有菜单，直接结束
	if len(menuIDs) == 0 {
		return nil
	}
	// 构建新关系
	roleMenus := make([]model.RoleMenu, 0, len(menuIDs))
	for _, v := range menuIDs {
		roleMenus = append(roleMenus, model.RoleMenu{
			RoleID: roleID,
			MenuID: v,
		})
	}
	// 批量插入
	if err := s.roleMenuRepo.CreateBatch(ctx, tx, roleMenus); err != nil {
		return err
	}
	return nil
}

// Delete 用于删除角色信息
func (s *Service) Delete(ctx context.Context, roleID uint64) (int, error) {
	// 获取角色信息
	role, err := s.roleRepo.GetByID(ctx, nil, roleID)
	if err != nil {
		return 60201, err
	}
	if role == nil {
		return 50201, nil
	}
	// 开启事务，删除角色
	err = s.db.Transaction(func(tx *gorm.DB) error {
		if err := s.roleMenuRepo.DeleteByRoleID(ctx, tx, roleID); err != nil {
			return err
		}
		if err := s.roleRepo.Delete(ctx, tx, roleID); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return 60203, err
	}
	// 返回数据
	return 0, nil
}

// RoleMenuTree 用于获取管理员权限内的菜单
func (s *Service) RoleMenuTree(ctx context.Context, roleID uint64, roleCode string) ([]RoleMenuItem, int, error) {
	// 获取菜单信息
	menus, err := s.getMenusByRole(ctx, roleID, roleCode)
	if err != nil {
		return nil, 60205, err
	}
	// 整理菜单信息
	list := make([]RoleMenuItem, 0, len(menus))
	for _, v := range menus {
		list = append(list, RoleMenuItem{
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

// getMenusByRole 用于获取指定角色拥有的菜单
func (s *Service) getMenusByRole(ctx context.Context, roleID uint64, roleCode string) ([]model.Menu, error) {
	// 超级管理员默认拥有全部菜单
	if roleCode == RoleCodeSuperAdmin {
		return s.menuRepo.GetEnableAll(ctx, nil)
	}
	// 非超管，根据角色ID获取菜单
	roleMenus, err := s.roleMenuRepo.GetByRoleID(ctx, nil, roleID)
	if err != nil {
		return nil, err
	}
	menuIDs := make([]uint64, 0, len(roleMenus))
	for _, v := range roleMenus {
		menuIDs = append(menuIDs, v.MenuID)
	}
	if len(menuIDs) == 0 {
		return nil, nil
	}
	return s.menuRepo.GetEnableByID(ctx, nil, menuIDs)
}

// buildTree 用于获取菜单权限树
func (s *Service) buildTree(list []RoleMenuItem, parentID uint64) []RoleMenuItem {
	// 构建 parent -> children map
	childrenMap := make(map[uint64][]RoleMenuItem)
	for _, el := range list {
		childrenMap[el.ParentID] = append(childrenMap[el.ParentID], el)
	}
	// 递归
	var build func(pid uint64) []RoleMenuItem
	build = func(pid uint64) []RoleMenuItem {
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
