package role

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

// ListPage 请求入参
type ListPageReq struct {
	PageResp
	Name   *string `json:"name"`
	Enable *int    `json:"enable"`
}

// ListPage 请求返回
type ListPageResp struct {
	Total    int64 `json:"total"`
	PageData []ListPageItem
}

type ListPageItem struct {
	ID     uint64 `json:"id"`
	Code   string `json:"code"`
	Name   string `json:"name"`
	Enable bool   `json:"enable"`
}

// Save 请求入参
type SaveReq struct {
	ID      *uint64  `json:"id"`
	Code    string   `json:"code"`
	Name    string   `json:"name"`
	Enable  bool     `json:"enable"`
	MenuIDs []uint64 `json:"menu_ids"`
}

// RoleMenuTreeByAdminID 请求返回
type RoleMenuTreeByAdminIDResp struct {
}

type RoleMenuItem struct {
	ID          uint64         `json:"id"`
	Code        string         `json:"code"`
	Enable      bool           `json:"enable"`
	Show        bool           `json:"show"`
	KeepAlive   bool           `json:"keep_alive"`
	Layout      string         `json:"layout"`
	Type        string         `json:"type"`
	ParentID    uint64         `json:"parent_id"`
	Name        string         `json:"name"`
	Icon        string         `json:"icon"`
	Path        string         `json:"path"`
	Component   string         `json:"component"`
	Order       int            `json:"order"`
	Redirect    string         `json:"redirect"`
	Method      string         `json:"method"`
	Description string         `json:"description"`
	Children    []RoleMenuItem `json:"children"`
}
