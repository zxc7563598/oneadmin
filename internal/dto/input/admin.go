package input

// AdminLoginReq 管理员登录请求
type AdminLoginReq struct {
	Username string  `json:"username" binding:"required" err:"required=10108"`
	Password string  `json:"password" binding:"required,min=6,max=32" err:"required=10109,min=10110,max=10111"`
	Captcha  *string `json:"captcha"`
}

// AdminRefreshReq 刷新管理员登录凭证请求
type AdminRefreshReq struct {
	Token string `json:"token" binding:"required" err:"required=10103"`
}

// AdminSwitchRoleReq 切换管理员角色请求
type AdminSwitchRoleReq struct {
	RoleID int `json:"role_id" binding:"required" err:"required=10112"`
}

// AdminChangePasswordReq 根据旧密码变更新密码请求
type AdminChangePasswordReq struct {
	OldPassword string `json:"old_password" binding:"required,min=6,max=32" err:"required=10109,min=10110,max=10111"`
	NewPassword string `json:"new_password" binding:"required,min=6,max=32" err:"required=10109,min=10110,max=10111"`
}

// AdminListPageReq 获取管理员分页列表请求
type AdminListPageReq struct {
	Page     int     `json:"page" binding:"required" err:"required=10103"`
	PageSize int     `json:"page_size" binding:"required" err:"required=10103"`
	Username *string `json:"username"`
	Gender   *int    `json:"gender"`
	Enable   *int    `json:"enable"`
}

// AdminSaveReq 变更管理员信息请求
type AdminSaveReq struct {
	ID       *uint64  `json:"id"`
	Enable   int      `json:"enable" binding:"required"  err:"required=10103"`
	Username string   `json:"username" binding:"required"  err:"required=10108"`
	Password *string  `json:"password"`
	RoleIds  []uint64 `json:"roleIds" binding:"required" err:"required=10112"`
}

// AdminDeleteReq 删除管理员信息请求
type AdminDeleteReq struct {
	ID uint64 `json:"id" binding:"required"`
}

// AdminResetAdminPasswordReq 变更管理员密码请求
type AdminResetAdminPasswordReq struct {
	ID       uint64 `json:"id" binding:"required"`
	Password string `json:"password" binding:"required,min=6,max=32" err:"required=10109,min=10110,max=10111"`
}

// AdminUpdateProfileReq 修改管理员个人信息请求
type AdminUpdateProfileReq struct {
	ID       uint64  `json:"id" binding:"required"`
	Nickname string  `json:"nickName" binding:"required" err:"required=10113"`
	Gender   int     `json:"gender" binding:"required"`
	Address  *string `json:"address"`
	Email    *string `json:"email"`
}
