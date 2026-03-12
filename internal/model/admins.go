package model

import "github.com/zxc7563598/oneadmin/internal/enum"

type Admin struct {
	ID       uint64      `gorm:"primaryKey"`
	Nickname string      `gorm:"type:varchar(100);not null;comment:昵称"`
	Username string      `gorm:"type:varchar(100);not null;comment:用户名"`
	Password string      `gorm:"type:varchar(255);not null;comment:密码"`
	Salt     int         `gorm:"type:int;not null;comment:扰乱码"`
	Token    *string     `gorm:"type:varchar(255);comment:登录凭证"`
	RoleID   *uint64     `gorm:"comment:当前角色ID"`
	Avatar   string      `gorm:"type:varchar(255);not null;default:'avatar.png';comment:头像"`
	Email    *string     `gorm:"type:varchar(255);comment:邮箱"`
	Address  *string     `gorm:"type:varchar(255);comment:地址"`
	Gender   enum.Gender `gorm:"type:smallint;not null;default:2;comment:性别"`
	Enable   enum.Enable `gorm:"type:smallint;not null;default:1;comment:是否启用"`
	BaseModel
}

func (Admin) TableName() string {
	return "admins"
}
