package input

// RoleListPageReq 获取角色分页列表请求
type RoleListPageReq struct {
	PageNo   int     `json:"pageNo" binding:"required" err:"required=10201"`
	PageSize int     `json:"pageSize" binding:"required" err:"required=10201"`
	Name     *string `json:"name"`
	Enable   *int    `json:"enable"`
}

// RoleSaveReq 处理创建或变更角色请求
type RoleSaveReq struct {
	ID            *uint64   `json:"id"`
	Code          *string   `json:"code"`
	Name          *string   `json:"name"`
	Enable        *bool     `json:"enable"`
	PermissionIds *[]uint64 `json:"permissionIds"`
}

// RoleDeleteReq 删除角色请求
type RoleDeleteReq struct {
	ID uint64 `json:"id" binding:"required" err:"required=10205"`
}
