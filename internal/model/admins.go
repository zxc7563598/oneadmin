package model

import "github.com/zxc7563598/oneadmin/internal/enum"

type Admin struct {
	ID       uint64      `gorm:"primaryKey"`
	Nickname string      `gorm:"type:varchar(100);not null;comment:昵称"`
	Username string      `gorm:"type:varchar(100);not null;comment:用户名"`
	Password string      `gorm:"type:varchar(255);not null;comment:密码"`
	Token    *string     `gorm:"type:varchar(255);comment:登录凭证"`
	RoleID   uint64      `gorm:"comment:当前角色ID"`
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

// AdminListQuery 用于后台列表查询入参，不对应数据库表
type AdminListQuery struct {
	Username *string
	Gender   *int
	Enable   *int
	Offset   int
	Limit    int
}

// AdminListItem 用于后台列表展示，不对应数据库表
type AdminListItem struct {
	ID        uint64
	Username  string
	Enable    enum.Enable
	Gender    enum.Gender
	Avatar    string
	Address   *string
	Email     *string
	RoleID    *uint64
	CreatedAt int64
	UpdatedAt int64
}

// AdminUpdateProfileForm 用于更新管理员个人资料，不对应数据库表
type AdminUpdateProfileForm struct {
	Nickname string
	Email    *string
	Address  *string
	Gender   enum.Gender
}
