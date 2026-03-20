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
	ListPage(ctx context.Context, tx *gorm.DB, query model.RoleListQuery) ([]model.RoleListItem, int64, error)
	UpdateByID(ctx context.Context, tx *gorm.DB, id uint64, queue model.RoleForm) error
}

// FindEnabled 获取全部启用数据
func (r *gormRepo) FindEnabled(ctx context.Context, tx *gorm.DB) ([]model.Role, error) {
	return r.FindByField(ctx, tx, "enable", enum.EnableEnable)
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
	updateMap := make(map[string]any)
	if queue.Code != nil {
		updateMap["code"] = *queue.Code
	}
	if queue.Name != nil {
		updateMap["name"] = *queue.Name
	}
	updateMap["enable"] = queue.Enable
	return r.UpdateMap(ctx, tx, "id", id, updateMap)
}
