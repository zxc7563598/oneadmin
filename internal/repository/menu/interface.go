package menu

import (
	"context"

	"github.com/zxc7563598/oneadmin/internal/enum"
	"github.com/zxc7563598/oneadmin/internal/model"
	"github.com/zxc7563598/oneadmin/internal/repository/base"
	"gorm.io/gorm"
)

type Repository interface {
	base.Repository[model.Menu]
	GetEnableAll(ctx context.Context, tx *gorm.DB) ([]model.Menu, error)
	GetEnableByID(ctx context.Context, tx *gorm.DB, ids []uint64) ([]model.Menu, error)
}

// GetEnableAll 获取全部启用菜单
func (r *gormRepo) GetEnableAll(ctx context.Context, tx *gorm.DB) ([]model.Menu, error) {
	return r.FindByField(ctx, tx, "enable", enum.EnableEnable)
}

// GetEnableByID 根据ID获取全部菜单
func (r *gormRepo) GetEnableByID(ctx context.Context, tx *gorm.DB, ids []uint64) ([]model.Menu, error) {
	db := r.getDB(ctx, tx)
	var list []model.Menu
	if err := db.Where("enable = ?", enum.EnableEnable).Where("id IN ?", ids).Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}
