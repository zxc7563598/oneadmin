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
}

// FindEnabled 获取全部启用数据
func (r *gormRepo) FindEnabled(ctx context.Context, tx *gorm.DB) ([]model.Role, error) {
	return r.FindByField(ctx, tx, "enable", enum.EnableEnable)
}
