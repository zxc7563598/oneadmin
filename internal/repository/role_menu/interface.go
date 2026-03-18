package role_menu

import (
	"context"

	"github.com/zxc7563598/oneadmin/internal/model"
	"github.com/zxc7563598/oneadmin/internal/repository/base"
	"gorm.io/gorm"
)

type Repository interface {
	base.Repository[model.RoleMenu]
	GetByRoleID(ctx context.Context, tx *gorm.DB, roleID uint64) ([]model.RoleMenu, error)
	DeleteByRoleID(ctx context.Context, tx *gorm.DB, roleID uint64) error
}

// GetByRoleID 根据角色ID获取权限信息
func (r *gormRepo) GetByRoleID(ctx context.Context, tx *gorm.DB, roleID uint64) ([]model.RoleMenu, error) {
	return r.FindByField(ctx, tx, "role_id", roleID)
}

// DeleteByRoleID 删除角色ID相关的权限
func (r *gormRepo) DeleteByRoleID(ctx context.Context, tx *gorm.DB, roleID uint64) error {
	db := r.getDB(ctx, tx)
	return db.WithContext(ctx).
		Where("role_id = ?", roleID).
		Delete(&model.RoleMenu{}).
		Error
}
