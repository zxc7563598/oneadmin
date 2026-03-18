package input

// RoleListPageReq 获取角色分页列表请求
type RoleListPageReq struct {
	Page     int     `json:"page" binding:"required" err:"required=10201"`
	PageSize int     `json:"page_size" binding:"required" err:"required=10201"`
	Name     *string `json:"name"`
	Enable   *int    `json:"enable"`
}

// RoleSaveReq 处理创建或变更角色请求
type RoleSaveReq struct {
	ID      *uint64  `json:"id"`
	Code    string   `json:"code" binding:"required" err:"required=10202"`
	Name    string   `json:"name" binding:"required" err:"required=10203"`
	Enable  bool     `json:"enable" binding:"required" err:"required=10204"`
	MenuIDs []uint64 `json:"menu_ids"`
}

// RoleDeleteReq 删除角色请求
type RoleDeleteReq struct {
	RoleID uint64 `json:"id" binding:"required" err:"required=10205"`
}
