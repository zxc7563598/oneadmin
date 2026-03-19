package role

import (
	"github.com/gin-gonic/gin"
	"github.com/zxc7563598/oneadmin/internal/dto/input"
	"github.com/zxc7563598/oneadmin/internal/dto/resp"
	"github.com/zxc7563598/oneadmin/internal/handler"
	"github.com/zxc7563598/oneadmin/internal/i18n"
	"github.com/zxc7563598/oneadmin/internal/logger"
	"github.com/zxc7563598/oneadmin/internal/response"
	"github.com/zxc7563598/oneadmin/internal/service/role"
	"go.uber.org/zap"
)

type Handler struct {
	roleSvc *role.Service
}

func New(roleSvc *role.Service) *Handler {
	return &Handler{
		roleSvc: roleSvc,
	}
}

// ListPage 处理获取角色列表(分页)请求
// POST /api/admin/roles/list
func (h *Handler) ListPage(c *gin.Context) {
	// 获取上下文/语言配置
	ctx := c.Request.Context()
	lang := i18n.GetLang(ctx)
	// 获取管理员ID
	adminInfo, ok := handler.GetAdminInfo(c)
	if !ok {
		response.Error(c, lang, 20001)
		return
	}
	// 获取请求参数
	var req input.RoleListPageReq
	if code, ok, err := handler.BindAndValidate(c, &req); !ok {
		handler.ErrorLog(
			logger.RoleLogger,
			"ListPage 参数异常",
			code,
			err,
		)
		response.Error(c, lang, code)
		return
	}
	// 执行请求
	svcResp, errCode, err := h.roleSvc.ListPage(ctx, role.ListPageReq{
		PageResp: role.PageResp{
			Page:     req.Page,
			PageSize: req.PageSize,
		},
		Name:   req.Name,
		Enable: req.Enable,
	})
	if errCode != 0 {
		handler.ErrorLog(
			logger.RoleLogger,
			"roleSvc.ListPage 调用失败",
			errCode,
			err,
			zap.Uint64("adminID", adminInfo.AdminID),
			zap.Int("req.page", req.Page),
			zap.Int("req.pageSize", req.PageSize),
			zap.Any("req.name", req.Name),
			zap.Any("req.enable", req.Enable),
		)
		response.Error(c, lang, errCode)
		return
	}
	// 返回结果
	response.Success(c, lang, resp.RoleListPageResp{
		Total:    svcResp.Total,
		PageData: toRoleListItems(svcResp.PageData),
	})
}

// ListPage 处理获取全部角色请求
// POST /api/admin/roles/all
func (h *Handler) ListAll(c *gin.Context) {
	// 获取上下文/语言配置
	ctx := c.Request.Context()
	lang := i18n.GetLang(ctx)
	// 获取管理员ID
	adminInfo, ok := handler.GetAdminInfo(c)
	if !ok {
		response.Error(c, lang, 20001)
		return
	}
	// 执行请求
	svcResp, errCode, err := h.roleSvc.ListAll(ctx)
	if errCode != 0 {
		handler.ErrorLog(
			logger.RoleLogger,
			"roleSvc.ListAll 调用失败",
			errCode,
			err,
			zap.Uint64("adminID", adminInfo.AdminID),
		)
		response.Error(c, lang, errCode)
		return
	}
	// 返回结果
	response.Success(c, lang, resp.RoleListAllResp{
		List: toRoleListItems(svcResp),
	})
}

// Save 处理创建或变更角色请求
// POST /api/admin/roles/save
func (h *Handler) Save(c *gin.Context) {
	// 获取上下文/语言配置
	ctx := c.Request.Context()
	lang := i18n.GetLang(ctx)
	// 获取管理员ID
	adminInfo, ok := handler.GetAdminInfo(c)
	if !ok {
		response.Error(c, lang, 20001)
		return
	}
	// 获取请求参数
	var req input.RoleSaveReq
	if code, ok, err := handler.BindAndValidate(c, &req); !ok {
		handler.ErrorLog(
			logger.RoleLogger,
			"Save 参数异常",
			code,
			err,
		)
		response.Error(c, lang, code)
		return
	}
	// 执行请求
	errCode, err := h.roleSvc.Save(ctx, role.SaveReq{
		ID:      req.ID,
		Code:    req.Code,
		Name:    req.Name,
		Enable:  req.Enable,
		MenuIDs: req.MenuIDs,
	})
	if errCode != 0 {
		handler.ErrorLog(
			logger.RoleLogger,
			"roleSvc.Save 调用失败",
			errCode,
			err,
			zap.Uint64("adminID", adminInfo.AdminID),
			zap.Any("req.id", req.ID),
			zap.String("req.code", req.Code),
			zap.String("req.name", req.Name),
			zap.Bool("req.enable", req.Enable),
			zap.Any("req.menu_ids", req.MenuIDs),
		)
		response.Error(c, lang, errCode)
		return
	}
	// 返回结果
	response.Success(c, lang, nil)
}

// Delete 删除角色请求
// POST /api/admin/roles/delete
func (h *Handler) Delete(c *gin.Context) {
	// 获取上下文/语言配置
	ctx := c.Request.Context()
	lang := i18n.GetLang(ctx)
	// 获取管理员ID
	adminInfo, ok := handler.GetAdminInfo(c)
	if !ok {
		response.Error(c, lang, 20001)
		return
	}
	// 获取请求参数
	var req input.RoleDeleteReq
	if code, ok, err := handler.BindAndValidate(c, &req); !ok {
		handler.ErrorLog(
			logger.RoleLogger,
			"Delete 参数异常",
			code,
			err,
		)
		response.Error(c, lang, code)
		return
	}
	// 执行请求
	errCode, err := h.roleSvc.Delete(ctx, req.RoleID)
	if errCode != 0 {
		handler.ErrorLog(
			logger.RoleLogger,
			"roleSvc.Delete 调用失败",
			errCode,
			err,
			zap.Uint64("adminID", adminInfo.AdminID),
			zap.Uint64("req.id", req.RoleID),
		)
		response.Error(c, lang, errCode)
		return
	}
	// 返回结果
	response.Success(c, lang, nil)
}

// Permissions 获取管理员菜单权限树请求
// POST /api/admin/roles/permissions
func (h *Handler) Permissions(c *gin.Context) {
	// 获取上下文/语言配置
	ctx := c.Request.Context()
	lang := i18n.GetLang(ctx)
	// 获取管理员ID
	adminInfo, ok := handler.GetAdminInfo(c)
	if !ok {
		response.Error(c, lang, 20001)
		return
	}
	// 获取角色权限内的菜单
	menus, errCode, err := h.roleSvc.RoleMenuTreeByAdminID(ctx, adminInfo.AdminID)
	if errCode != 0 {
		handler.ErrorLog(
			logger.RoleLogger,
			"roleSvc.RoleMenuTreeByAdminID 调用失败",
			errCode,
			err,
			zap.Uint64("adminID", adminInfo.AdminID),
		)
		response.Error(c, lang, errCode)
		return
	}
	// 返回结果
	response.Success(c, lang, resp.RolePermissionsResp{
		Menu: toRoleMenuItem(menus),
	})
}
