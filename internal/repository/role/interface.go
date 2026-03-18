package role

import (
	"context"

	"github.com/zxc7563598/oneadmin/internal/enum"
	"github.com/zxc7563598/oneadmin/internal/model"
	"github.com/zxc7563598/oneadmin/internal/repository/base"
	"gorm.io/gorm"
)

type Repository interface {
	base.Repository[model.Role]
	FindEnabled(ctx context.Context, tx *gorm.DB) ([]model.Role, error)
	RoleByAdminID(ctx context.Context, tx *gorm.DB, adminID uint64) (*model.Role, error)
	ListPage(ctx context.Context, tx *gorm.DB, query model.RoleListQuery) ([]model.RoleListItem, int64, error)
	UpdateByID(ctx context.Context, tx *gorm.DB, id uint64, queue model.RoleForm) error
}

// FindEnabled 获取全部启用数据
func (r *gormRepo) FindEnabled(ctx context.Context, tx *gorm.DB) ([]model.Role, error) {
	return r.FindByField(ctx, tx, "enable", enum.EnableEnable)
}

// RoleByAdminID 根据管理员ID获取角色信息
func (r *gormRepo) RoleByAdminID(ctx context.Context, tx *gorm.DB, adminID uint64) (*model.Role, error) {
	return r.FindOneByField(ctx, tx, "admin_id", adminID)
}

// ListPage 获取分页列表数据
func (r *gormRepo) ListPage(ctx context.Context, tx *gorm.DB, query model.RoleListQuery) ([]model.RoleListItem, int64, error) {
	var list []model.RoleListItem
	var total int64
	db := r.getDB(ctx, tx)
	db = db.Model(&model.Role{})
	if query.Name != nil && *query.Name != "" {
		db = db.Where("name LIKE ?", "%"+*query.Name+"%")
	}
	if query.Enable != nil {
		db = db.Where("enable = ?", *query.Enable)
	}
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err := db.Order("id asc").Offset(query.Offset).Limit(query.Limit).Find(&list).Error
	if err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

// UpdateByID 变更基本信息
func (r *gormRepo) UpdateByID(ctx context.Context, tx *gorm.DB, id uint64, queue model.RoleForm) error {
	return r.UpdateMap(ctx, tx, "id", id, map[string]any{
		"code":   queue.Code,
		"name":   queue.Name,
		"enable": queue.Enable,
	})
}
