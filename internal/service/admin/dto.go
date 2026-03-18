package admin

// 通用分页请求参数
type PageResp struct {
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
}

func (r *PageResp) OffsetLimit() (int, int) {
	if r.Page < 1 {
		r.Page = 1
	}
	if r.PageSize < 1 {
		r.PageSize = 10
	}
	if r.PageSize > 100 {
		r.PageSize = 100
	}
	offset := (r.Page - 1) * r.PageSize
	return offset, r.PageSize
}

// Login 请求返回
type LoginResp struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// RefreshLogin 请求返回
type RefreshLoginResp struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// ListPage 请求入参
type ListPageReq struct {
	PageResp
	Username *string `json:"username"`
	Gender   *int    `json:"gender"`
	Enable   *int    `json:"enable"`
}

// ListPage 请求返回
type ListPageResp struct {
	Total    int `json:"total"`
	PageData []ListPageItem
}

type ListPageItem struct {
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

// Details 请求返回
type DetailsResp struct {
	ID          uint64             `json:"id"`
	Username    string             `json:"username"`
	Enable      bool               `json:"enable"`
	CreatedAt   string             `json:"created_at"`
	UpdatedAt   string             `json:"updated_at"`
	Profile     DetailsProfileItem `json:"profile"`
	Roles       []DetailsRoleItem  `json:"roles"`
	CurrentRole DetailsRoleItem    `json:"currentRole"`
}

type DetailsRoleItem struct {
	ID     uint64 `json:"id"`
	Code   string `json:"code"`
	Name   string `json:"name"`
	Enable string `json:"enable"`
}

type DetailsProfileItem struct {
	ID       uint64 `json:"id"`
	Nickname string `json:"nickName"`
	Gender   int    `json:"gender"`
	Avatar   string `json:"avatar"`
	Address  string `json:"address"`
	Email    string `json:"email"`
}

// Save 请求入参
type SaveReq struct {
	ID       *uint64  `json:"id"`
	Enable   int      `json:"enable"`
	Username string   `json:"username"`
	Password *string  `json:"password"`
	RoleIds  []uint64 `json:"roleIds"`
}

// UpdateProfile 请求入参
type UpdateProfileReq struct {
	ID       uint64  `json:"id"`
	Nickname string  `json:"nickName"`
	Gender   int     `json:"gender"`
	Address  *string `json:"address"`
	Email    *string `json:"email"`
}
