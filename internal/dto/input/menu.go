package input

// 验证菜单是否存在请求
type MenuValidateReq struct {
	Path string `json:"path" binding:"required,min=1" err:"required=10301,min=10302"`
}

// 获取菜单下的按钮请求
type MenuButtonsReq struct {
	ParentID uint64 `json:"parent_id" binding:"required" err:"required=10301"`
}

// 添加或变更菜单信息请求
type MenuSaveReq struct {
	ID        *uint64 `json:"id"`
	Code      string  `json:"code" binding:"required,min=1" err:"required=10301,min=10303"`
	Enable    bool    `json:"enable" binding:"required" err:"required=10301"`
	Show      bool    `json:"show" binding:"required" err:"required=10301"`
	KeepAlive bool    `json:"keep_alive" binding:"required" err:"required=10301"`
	Layout    string  `json:"layout" binding:"required" err:"required=10301"`
	Type      string  `json:"type" binding:"required,min=1" err:"required=10301,min=10304"`
	ParentID  uint64  `json:"parent_id"`
	Name      string  `json:"name" binding:"required,min=1" err:"required=10301,min=10305"`
	Icon      string  `json:"icon"`
	Path      string  `json:"path"`
	Component string  `json:"component"`
	Order     int     `json:"order"`
}

// 快速切换菜单的启用状态请求
type MenuToggleReq struct {
	ID uint64 `json:"id" binding:"required" err:"required=10301"`
}

// 删除菜单请求
type MenuDeleteReq struct {
	ID uint64 `json:"id" binding:"required" err:"required=10301"`
}
