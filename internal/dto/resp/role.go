package resp

// RoleListPageResp 获取管理员分页列表返回
type RoleListPageResp struct {
	Total    int64              `json:"total"`
	PageData []RoleListPageItem `json:"pageData"`
}

type RoleListPageItem struct {
	ID            uint64   `json:"id"`
	Code          string   `json:"code"`
	Name          string   `json:"name"`
	Enable        bool     `json:"enable"`
	PermissionIds []uint64 `json:"permissionIds"`
}

// RoleListAllResp 获取全部角色请求返回
type RoleListAllResp struct {
	List []RoleListPageItem `json:"list"`
}

// RolePermissionsResp 获取管理员菜单权限树请求返回
type RolePermissionsResp struct {
	Menu []RoleMenuItem `json:"menu"`
}

type RoleMenuItem struct {
	ID          uint64         `json:"id"`
	Code        string         `json:"code"`
	Enable      bool           `json:"enable"`
	Show        bool           `json:"show"`
	KeepAlive   bool           `json:"keepAlive"`
	Layout      string         `json:"layout"`
	Type        string         `json:"type"`
	ParentID    uint64         `json:"parentId"`
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
