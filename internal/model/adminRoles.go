package model

type AdminRole struct {
	ID      uint64 `gorm:"primaryKey"`
	AdminID uint64 `gorm:"not null;default:0;comment:管理员ID"`
	RoleID  uint64 `gorm:"not null;default:0;comment:角色ID"`
	BaseModel
}

func (AdminRole) TableName() string {
	return "admin_roles"
}
