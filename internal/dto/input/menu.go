package input

// 验证菜单是否存在请求
type MenuValidateReq struct {
	Path string `json:"path"`
}

// 获取菜单下的按钮请求
type MenuButtonsReq struct {
	ParentID uint64 `json:"parent_id"`
}

// 添加或变更菜单信息请求
type MenuSaveReq struct {
	ID        *uint64 `json:"id"`
	Code      string  `json:"code"`
	Enable    bool    `json:"enable"`
	Show      bool    `json:"show"`
	KeepAlive bool    `json:"keep_alive"`
	Layout    string  `json:"layout"`
	Type      string  `json:"type"`
	ParentID  uint64  `json:"parent_id"`
	Name      string  `json:"name"`
	Icon      string  `json:"icon"`
	Path      string  `json:"path"`
	Component string  `json:"component"`
	Order     int     `json:"order"`
}

// 快速切换菜单的启用状态请求
type MenuToggleReq struct {
	ID uint64 `json:"id"`
}

// 删除菜单请求
type MenuDeleteReq struct {
	ID uint64 `json:"id"`
}
