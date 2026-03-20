package admin

import (
	"context"

	"github.com/zxc7563598/oneadmin/internal/enum"
	"github.com/zxc7563598/oneadmin/internal/model"
	"github.com/zxc7563598/oneadmin/internal/repository/base"
	"gorm.io/gorm"
)

type Repository interface {
	base.Repository[model.Admin]
	GetByUsername(ctx context.Context, tx *gorm.DB, username string) (*model.Admin, error)
	UpdateToken(ctx context.Context, tx *gorm.DB, id uint64, token *string) error
	UpdateRole(ctx context.Context, tx *gorm.DB, adminID, roleID uint64) error
	UpdatePassword(ctx context.Context, tx *gorm.DB, adminID uint64, password string) error
	UpdateInfo(ctx context.Context, tx *gorm.DB, adminID uint64, Username string, Enable enum.Enable) error
	ListPage(ctx context.Context, tx *gorm.DB, query model.AdminListQuery) ([]model.AdminListItem, int64, error)
	UpdateProfile(ctx context.Context, tx *gorm.DB, adminID uint64, form model.AdminUpdateProfileForm) error
}

// GetByUsername 根据 username 获取账号信息
func (r *gormRepo) GetByUsername(ctx context.Context, tx *gorm.DB, username string) (*model.Admin, error) {
	return r.FindOneByField(ctx, tx, "username", username)
}

// UpdateToken 根据 id 更换管理员 refreshToken
func (r *gormRepo) UpdateToken(ctx context.Context, tx *gorm.DB, id uint64, token *string) error {
	return r.UpdateField(ctx, tx, id, "token", token)
}

// UpdateRole 设置管理员的角色
func (r *gormRepo) UpdateRole(ctx context.Context, tx *gorm.DB, adminID, roleID uint64) error {
	return r.UpdateField(ctx, tx, adminID, "role_id", roleID)
}

// UpdatePassword 更新管理员密码
func (r *gormRepo) UpdatePassword(ctx context.Context, tx *gorm.DB, adminID uint64, password string) error {
	return r.UpdateField(ctx, tx, adminID, "password", password)
}

// UpdateInfo 更新管理员基本信息
func (r *gormRepo) UpdateInfo(ctx context.Context, tx *gorm.DB, adminID uint64, Username string, Enable enum.Enable) error {
	return r.UpdateMap(ctx, tx, "id", adminID, map[string]any{
		"enable":   Enable,
		"username": Username,
	})
}

// ListPage 获取分页列表数据
func (r *gormRepo) ListPage(ctx context.Context, tx *gorm.DB, query model.AdminListQuery) ([]model.AdminListItem, int64, error) {
	var list []model.AdminListItem
	var total int64
	db := r.getDB(ctx, tx)
	db = db.Model(&model.Admin{})
	if query.Username != nil && *query.Username != "" {
		db = db.Where("username LIKE ?", "%"+*query.Username+"%")
	}
	if query.Gender != nil {
		db = db.Where("gender = ?", *query.Gender)
	}
	if query.Enable != nil {
		db = db.Where("enable = ?", *query.Enable)
	}
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err := db.Order("created_at desc").Offset(query.Offset).Limit(query.Limit).Find(&list).Error
	if err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

// UpdateProfile 更新管理员个人资料
func (r *gormRepo) UpdateProfile(ctx context.Context, tx *gorm.DB, adminID uint64, form model.AdminUpdateProfileForm) error {
	updateMap := make(map[string]any)
	if form.Nickname != nil {
		updateMap["nickname"] = *form.Nickname
	}
	if form.Email != nil {
		updateMap["email"] = *form.Email
	}
	if form.Address != nil {
		updateMap["address"] = *form.Address
	}
	if form.Gender != nil {
		updateMap["gender"] = enum.Enable(*form.Gender)
	}
	if form.Avatar != nil {
		updateMap["avatar"] = *form.Avatar
	}
	return r.UpdateMap(ctx, tx, "id", adminID, updateMap)
}
