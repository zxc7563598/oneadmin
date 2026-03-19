package menu

import (
	"github.com/gin-gonic/gin"
	"github.com/zxc7563598/oneadmin/internal/dto/input"
	"github.com/zxc7563598/oneadmin/internal/dto/resp"
	"github.com/zxc7563598/oneadmin/internal/handler"
	"github.com/zxc7563598/oneadmin/internal/i18n"
	"github.com/zxc7563598/oneadmin/internal/logger"
	"github.com/zxc7563598/oneadmin/internal/response"
	"github.com/zxc7563598/oneadmin/internal/service/menu"
	"go.uber.org/zap"
)

type Handler struct {
	menuSvc *menu.Service
}

func New(menuSvc *menu.Service) *Handler {
	return &Handler{
		menuSvc: menuSvc,
	}
}

// List 处理获取全部菜单请求
// POST /api/admin/menu/list
func (h *Handler) List(c *gin.Context) {
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
	svcResp, errCode, err := h.menuSvc.MenuTree(ctx)
	if errCode != 0 {
		handler.ErrorLog(
			logger.MenuLogger,
			"menuSvc.MenuTree 调用失败",
			errCode,
			err,
			zap.Uint64("adminID", adminInfo.AdminID),
		)
		response.Error(c, lang, errCode)
		return
	}
	// 返回结果
	response.Success(c, lang, resp.MenuListResp{
		Menu: toMenuItem(svcResp),
	})
}

// Validate 处理验证菜单是否存在
// POST /api/admin/menu/validate
func (h *Handler) Validate(c *gin.Context) {
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
	var req input.MenuValidateReq
	if code, ok, err := handler.BindAndValidate(c, &req); !ok {
		handler.ErrorLog(
			logger.MenuLogger,
			"Validate 参数异常",
			code,
			err,
		)
		response.Error(c, lang, code)
		return
	}
	// 执行请求
	svcResp, errCode, err := h.menuSvc.MenuExists(ctx, req.Path)
	if errCode != 0 {
		handler.ErrorLog(
			logger.MenuLogger,
			"menuSvc.MenuExists 调用失败",
			errCode,
			err,
			zap.Uint64("adminID", adminInfo.AdminID),
			zap.String("req.path", req.Path),
		)
		response.Error(c, lang, errCode)
		return
	}
	// 返回结果
	response.Success(c, lang, resp.MenuValidateResp{
		Has: svcResp,
	})
}

// Buttons 获取菜单下的按钮
// POST /api/admin/menu/buttons
func (h *Handler) Buttons(c *gin.Context) {
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
	var req input.MenuButtonsReq
	if code, ok, err := handler.BindAndValidate(c, &req); !ok {
		handler.ErrorLog(
			logger.MenuLogger,
			"Buttons 参数异常",
			code,
			err,
		)
		response.Error(c, lang, code)
		return
	}
	// 执行请求
	svcResp, errCode, err := h.menuSvc.MenuButtons(ctx, req.ParentID)
	if errCode != 0 {
		handler.ErrorLog(
			logger.MenuLogger,
			"menuSvc.MenuButtons 调用失败",
			errCode,
			err,
			zap.Uint64("adminID", adminInfo.AdminID),
			zap.Uint64("req.parent_id", req.ParentID),
		)
		response.Error(c, lang, errCode)
		return
	}
	// 返回结果
	response.Success(c, lang, resp.MenuButtonsResp{
		Menu: toMenuItem(svcResp),
	})
}

// Save 添加或变更菜单信息
// POST /api/admin/menu/save
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
	var req input.MenuSaveReq
	if code, ok, err := handler.BindAndValidate(c, &req); !ok {
		handler.ErrorLog(
			logger.MenuLogger,
			"Buttons 参数异常",
			code,
			err,
		)
		response.Error(c, lang, code)
		return
	}
	// 执行请求
	errCode, err := h.menuSvc.Save(ctx, menu.SaveReq{
		ID:        req.ID,
		Code:      req.Code,
		Enable:    req.Enable,
		Show:      req.Show,
		KeepAlive: req.KeepAlive,
		Layout:    req.Layout,
		Type:      req.Type,
		ParentID:  req.ParentID,
		Name:      req.Name,
		Icon:      req.Icon,
		Path:      req.Path,
		Component: req.Component,
		Order:     req.Order,
	})
	if errCode != 0 {
		handler.ErrorLog(
			logger.MenuLogger,
			"menuSvc.Save 调用失败",
			errCode,
			err,
			zap.Uint64("adminID", adminInfo.AdminID),
			zap.Any("req.id", req.ID),
			zap.String("req.code", req.Code),
			zap.Bool("req.enable", req.Enable),
			zap.Bool("req.show", req.Show),
			zap.Bool("req.keep_alive", req.KeepAlive),
			zap.String("req.layout", req.Layout),
			zap.String("req.type", req.Type),
			zap.Uint64("req.parent_id", req.ParentID),
			zap.String("req.name", req.Name),
			zap.String("req.icon", req.Icon),
			zap.String("req.path", req.Path),
			zap.String("req.component", req.Component),
			zap.Int("req.order", req.Order),
		)
		response.Error(c, lang, errCode)
		return
	}
	// 返回结果
	response.Success(c, lang, nil)
}

// Toggle 快速切换菜单的启用状态
// POST /api/admin/menu/toggle
func (h *Handler) Toggle(c *gin.Context) {
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
	var req input.MenuToggleReq
	if code, ok, err := handler.BindAndValidate(c, &req); !ok {
		handler.ErrorLog(
			logger.MenuLogger,
			"Buttons 参数异常",
			code,
			err,
		)
		response.Error(c, lang, code)
		return
	}
	// 执行请求
	errCode, err := h.menuSvc.SetMenuEnable(ctx, req.ID)
	if errCode != 0 {
		handler.ErrorLog(
			logger.MenuLogger,
			"menuSvc.SetMenuEnable 调用失败",
			errCode,
			err,
			zap.Uint64("adminID", adminInfo.AdminID),
			zap.Uint64("req.id", req.ID),
		)
		response.Error(c, lang, errCode)
		return
	}
	// 返回结果
	response.Success(c, lang, nil)
}

// Delete 删除菜单
// POST /api/admin/menu/delete
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
	var req input.MenuDeleteReq
	if code, ok, err := handler.BindAndValidate(c, &req); !ok {
		handler.ErrorLog(
			logger.MenuLogger,
			"Buttons 参数异常",
			code,
			err,
		)
		response.Error(c, lang, code)
		return
	}
	// 执行请求
	errCode, err := h.menuSvc.Delete(ctx, req.ID)
	if errCode != 0 {
		handler.ErrorLog(
			logger.MenuLogger,
			"menuSvc.Delete 调用失败",
			errCode,
			err,
			zap.Uint64("adminID", adminInfo.AdminID),
			zap.Uint64("req.id", req.ID),
		)
		response.Error(c, lang, errCode)
		return
	}
	// 返回结果
	response.Success(c, lang, nil)
}
