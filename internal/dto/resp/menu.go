package resp

// MenuItem 单条菜单信息
type MenuItem struct {
	ID          uint64     `json:"id"`
	Code        string     `json:"code"`
	Enable      bool       `json:"enable"`
	Show        bool       `json:"show"`
	KeepAlive   bool       `json:"keepAlive"`
	Layout      string     `json:"layout"`
	Type        string     `json:"type"`
	ParentID    uint64     `json:"parentId"`
	Name        string     `json:"name"`
	Icon        string     `json:"icon"`
	Path        string     `json:"path"`
	Component   string     `json:"component"`
	Order       int        `json:"order"`
	Redirect    string     `json:"redirect"`
	Method      string     `json:"method"`
	Description string     `json:"description"`
	Children    []MenuItem `json:"children"`
}

// MenuListResp 获取全部菜单返回
type MenuListResp struct {
	Menu []MenuItem `json:"menu"`
}

// MenuValidateResp 验证菜单是否存在返回
type MenuValidateResp struct {
	Has bool `json:"has"`
}

// MenuButtonsResp
type MenuButtonsResp struct {
	Total    int64      `json:"total"`
	PageData []MenuItem `json:"pageData"`
}
