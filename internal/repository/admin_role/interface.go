package admin_role

import (
	"context"

	"github.com/zxc7563598/oneadmin/internal/model"
	"github.com/zxc7563598/oneadmin/internal/repository/base"
	"gorm.io/gorm"
)

type Repository interface {
	base.Repository[model.AdminRole]
	HasRole(ctx context.Context, tx *gorm.DB, adminID, roleID uint64) (bool, error)
	DeleteByAdminID(ctx context.Context, tx *gorm.DB, adminID uint64) error
	GetByRoleIDs(ctx context.Context, tx *gorm.DB, ids []uint64) ([]model.AdminRole, error)
}

// HasRole 判断管理员是否拥有指定角色
func (r *gormRepo) HasRole(ctx context.Context, tx *gorm.DB, adminID, roleID uint64) (bool, error) {
	db := r.getDB(ctx, tx)
	var count int64
	if err := db.Model(&model.AdminRole{}).Where("admin_id = ? AND role_id = ?", adminID, roleID).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// DeleteByAdminID 删除管理员ID相关的角色
func (r *gormRepo) DeleteByAdminID(ctx context.Context, tx *gorm.DB, adminID uint64) error {
	db := r.getDB(ctx, tx)
	return db.WithContext(ctx).
		Where("admin_id = ?", adminID).
		Delete(&model.AdminRole{}).
		Error
}

// GetByAdminIDs 根据管理员ID批量获取
func (r *gormRepo) GetByRoleIDs(ctx context.Context, tx *gorm.DB, ids []uint64) ([]model.AdminRole, error) {
	db := r.getDB(ctx, tx)
	var list []model.AdminRole
	if err := db.Where("admin_id IN ?", ids).Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}
