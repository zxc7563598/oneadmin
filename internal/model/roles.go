package model

import "github.com/zxc7563598/oneadmin/internal/enum"

type Role struct {
	ID     uint64      `gorm:"primaryKey"`
	Code   string      `gorm:"type:varchar(100);not null;comment:编码"`
	Name   string      `gorm:"type:varchar(100);not null;comment:名称"`
	Enable enum.Enable `gorm:"type:smallint;not null;default:0;comment:是否启用"`
	BaseModel
}

func (Role) TableName() string {
	return "roles"
}
