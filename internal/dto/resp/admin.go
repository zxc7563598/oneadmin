package resp

// AdminLoginResp 管理员登录/刷新登录返回
type AdminLoginResp struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// AdminListPageResp 获取管理员分页列表返回
type AdminListPageResp struct {
	Total    int64 `json:"total"`
	PageData []AdminListPageItem
}

type AdminListPageItem struct {
	ID        uint64  `json:"id"`
	Username  string  `json:"username"`
	Enable    string  `json:"enable"`
	Gender    string  `json:"gender"`
	Avatar    string  `json:"avatar"`
	Address   *string `json:"address"`
	Email     *string `json:"email"`
	RoleID    *uint64 `json:"role_id"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
}

// AdminDetailsResp 获取单个管理员详细信息返回
type AdminDetailsResp struct {
	ID          uint64                  `json:"id"`
	Username    string                  `json:"username"`
	Enable      bool                    `json:"enable"`
	CreatedAt   string                  `json:"created_at"`
	UpdatedAt   string                  `json:"updated_at"`
	Profile     AdminDetailsProfileItem `json:"profile"`
	Roles       []AdminDetailsRoleItem  `json:"roles"`
	CurrentRole AdminDetailsRoleItem    `json:"currentRole"`
}

type AdminDetailsRoleItem struct {
	ID     uint64 `json:"id"`
	Code   string `json:"code"`
	Name   string `json:"name"`
	Enable string `json:"enable"`
}

type AdminDetailsProfileItem struct {
	ID       uint64 `json:"id"`
	Nickname string `json:"nickName"`
	Gender   int    `json:"gender"`
	Avatar   string `json:"avatar"`
	Address  string `json:"address"`
	Email    string `json:"email"`
}
