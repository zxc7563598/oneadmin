package input

// AdminLoginReq 管理员登录请求
type AdminLoginReq struct {
	Username string  `json:"username" binding:"required,min=1" err:"required=10102,min=10102"`
	Password string  `json:"password" binding:"required,min=6,max=32" err:"required=10103,min=10104,max=10105"`
	Captcha  *string `json:"captcha"`
}

// AdminRefreshReq 刷新管理员登录凭证请求
type AdminRefreshReq struct {
	Token string `json:"token" binding:"required" err:"required=10101"`
}

// AdminSwitchRoleReq 切换管理员角色请求
type AdminSwitchRoleReq struct {
	RoleID int `json:"roleID" binding:"required" err:"required=10106"`
}

// AdminChangePasswordReq 根据旧密码变更新密码请求
type AdminChangePasswordReq struct {
	OldPassword string `json:"oldPassword" binding:"required,min=6,max=32" err:"required=10103,min=10104,max=10105"`
	NewPassword string `json:"newPassword" binding:"required,min=6,max=32" err:"required=10103,min=10104,max=10105"`
}

// AdminListPageReq 获取管理员分页列表请求
type AdminListPageReq struct {
	PageNo   int     `json:"pageNo" binding:"required" err:"required=10101"`
	PageSize int     `json:"pageSize" binding:"required" err:"required=10101"`
	Username *string `json:"username"`
	Gender   *int    `json:"gender"`
	Enable   *int    `json:"enable"`
}

// AdminSaveReq 变更管理员信息请求
type AdminSaveReq struct {
	ID       *uint64  `json:"id"`
	Enable   *int     `json:"enable"`
	Username string   `json:"username" binding:"required,min=1" err:"required=10102,min=10102"`
	Password *string  `json:"password"`
	RoleIds  []uint64 `json:"roleIds" binding:"required,min=1" err:"required=10106,min=10108"`
}

// AdminDeleteReq 删除管理员信息请求
type AdminDeleteReq struct {
	ID uint64 `json:"id" binding:"required" err:"required=10101"`
}

// AdminResetAdminPasswordReq 变更管理员密码请求
type AdminResetAdminPasswordReq struct {
	ID       uint64 `json:"id" binding:"required"`
	Password string `json:"password" binding:"required,min=6,max=32" err:"required=10103,min=10104,max=10105"`
}

// AdminUpdateProfileReq 修改管理员个人信息请求
type AdminUpdateProfileReq struct {
	ID       uint64  `json:"id" binding:"required" err:"required=10101"`
	Nickname *string `json:"nickName"`
	Gender   *int    `json:"gender"`
	Address  *string `json:"address"`
	Email    *string `json:"email"`
	Avatar   *string `json:"avatar"`
}
