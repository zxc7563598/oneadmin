package admin

import (
	"github.com/gin-gonic/gin"
	"github.com/zxc7563598/oneadmin/internal/dto/input"
	"github.com/zxc7563598/oneadmin/internal/dto/resp"
	"github.com/zxc7563598/oneadmin/internal/handler"
	"github.com/zxc7563598/oneadmin/internal/i18n"
	"github.com/zxc7563598/oneadmin/internal/logger"
	"github.com/zxc7563598/oneadmin/internal/response"
	"github.com/zxc7563598/oneadmin/internal/service/admin"
	"go.uber.org/zap"
)

type Handler struct {
	adminSvc admin.Service
}

func New(adminSvc *admin.Service) *Handler {
	return &Handler{
		adminSvc: *adminSvc,
	}
}

// Login 处理管理员登录请求
// POST /api/admin/auth/login
func (h *Handler) Login(c *gin.Context) {
	ctx := c.Request.Context()
	lang := i18n.GetLang(ctx)
	// 获取请求参数
	var req input.AdminLoginReq
	if code, ok, err := handler.BindAndValidate(c, &req); !ok {
		handler.ErrorLog(
			logger.AdminLogger,
			"Login 参数异常",
			code,
			err,
		)
		response.Error(c, lang, code)
		return
	}
	// 执行请求
	svcResp, errCode, err := h.adminSvc.Login(ctx, req.Username, req.Password, req.Captcha)
	if errCode != 0 {
		handler.ErrorLog(
			logger.AdminLogger,
			"adminSvc.Login 调用失败",
			errCode,
			err,
			zap.String("uname", req.Username),
		)
		response.Error(c, lang, errCode)
		return
	}
	// 返回数据
	response.Success(c, lang, resp.AdminLoginResp{
		AccessToken:  svcResp.AccessToken,
		RefreshToken: svcResp.RefreshToken,
	})
}

// Refresh 处理管理员刷新登陆凭证请求
// POST /api/admin/auth/refresh
func (h *Handler) Refresh(c *gin.Context) {
	// 获取上下文/语言配置
	ctx := c.Request.Context()
	lang := i18n.GetLang(ctx)
	// 获取请求参数
	var req input.AdminRefreshReq
	if code, ok, err := handler.BindAndValidate(c, &req); !ok {
		handler.ErrorLog(
			logger.AdminLogger,
			"Refresh 参数异常",
			code,
			err,
		)
		response.Error(c, lang, code)
		return
	}
	// 执行请求
	svcResp, errCode, err := h.adminSvc.RefreshLogin(ctx, req.Token)
	if errCode != 0 {
		handler.ErrorLog(
			logger.AdminLogger,
			"adminSvc.RefreshLogin 调用失败",
			errCode,
			err,
		)
		response.Error(c, lang, errCode)
		return
	}
	// 返回结果
	response.Success(c, lang, resp.AdminLoginResp{
		AccessToken:  svcResp.AccessToken,
		RefreshToken: svcResp.RefreshToken,
	})
}

// Logout 处理管理员退出登录请求
// POST /api/admin/auth/logout
func (h *Handler) Logout(c *gin.Context) {
	// 获取上下文/语言配置
	ctx := c.Request.Context()
	lang := i18n.GetLang(ctx)
	// 获取管理员ID
	adminID, ok := handler.GetAdminID(c)
	if !ok {
		response.Error(c, lang, 20101)
		return
	}
	// 执行请求
	errCode, err := h.adminSvc.Logout(ctx, adminID)
	if errCode != 0 {
		handler.ErrorLog(
			logger.AdminLogger,
			"adminSvc.Logout 调用失败",
			errCode,
			err,
			zap.Uint64("adminID", adminID),
		)
		response.Error(c, lang, errCode)
		return
	}
	// 返回结果
	response.Success(c, lang, nil)
}

// SwitchRole 处理管理员切换角色请求
// POST /api/admin/auth/switch-role
func (h *Handler) SwitchRole(c *gin.Context) {
	// 获取上下文/语言配置
	ctx := c.Request.Context()
	lang := i18n.GetLang(ctx)
	// 获取管理员ID
	adminID, ok := handler.GetAdminID(c)
	if !ok {
		response.Error(c, lang, 20101)
		return
	}
	// 获取请求参数
	var req input.AdminSwitchRoleReq
	if code, ok, err := handler.BindAndValidate(c, &req); !ok {
		handler.ErrorLog(
			logger.AdminLogger,
			"SwitchRole 参数异常",
			code,
			err,
		)
		response.Error(c, lang, code)
		return
	}
	// 执行请求
	errCode, err := h.adminSvc.SwitchRole(ctx, adminID, uint64(req.RoleID))
	if errCode != 0 {
		handler.ErrorLog(
			logger.AdminLogger,
			"adminSvc.SwitchRole 调用失败",
			errCode,
			err,
			zap.Uint64("adminID", adminID),
			zap.Uint64("RoleID", uint64(req.RoleID)),
		)
		response.Error(c, lang, errCode)
		return
	}
	// 返回结果
	response.Success(c, lang, nil)
}

// ChangePassword 处理管理员根据旧密码变更新密码请求
// POST /api/admin/auth/change-password
func (h *Handler) ChangePassword(c *gin.Context) {
	// 获取上下文/语言配置
	ctx := c.Request.Context()
	lang := i18n.GetLang(ctx)
	// 获取管理员ID
	adminID, ok := handler.GetAdminID(c)
	if !ok {
		response.Error(c, lang, 20101)
		return
	}
	// 获取请求参数
	var req input.AdminChangePasswordReq
	if code, ok, err := handler.BindAndValidate(c, &req); !ok {
		handler.ErrorLog(
			logger.AdminLogger,
			"ChangePassword 参数异常",
			code,
			err,
		)
		response.Error(c, lang, code)
		return
	}
	// 执行请求
	errCode, err := h.adminSvc.ChangePassword(ctx, adminID, req.OldPassword, req.NewPassword)
	if errCode != 0 {
		handler.ErrorLog(
			logger.AdminLogger,
			"adminSvc.ChangePassword 调用失败",
			errCode,
			err,
			zap.Uint64("adminID", adminID),
		)
		response.Error(c, lang, errCode)
		return
	}
	// 返回结果
	response.Success(c, lang, nil)
}

// ListPage 处理获取管理员列表(分页)请求
// POST /api/admin/admin/list
func (h *Handler) ListPage(c *gin.Context) {
	// 获取上下文/语言配置
	ctx := c.Request.Context()
	lang := i18n.GetLang(ctx)
	// 获取管理员ID
	adminID, ok := handler.GetAdminID(c)
	if !ok {
		response.Error(c, lang, 20101)
		return
	}
	// 获取请求参数
	var req input.AdminListPageReq
	if code, ok, err := handler.BindAndValidate(c, &req); !ok {
		handler.ErrorLog(
			logger.AdminLogger,
			"ListPage 参数异常",
			code,
			err,
		)
		response.Error(c, lang, code)
		return
	}
	// 执行请求
	svcResp, errCode, err := h.adminSvc.ListPage(ctx, admin.ListPageReq{
		PageResp: admin.PageResp{
			Page:     req.Page,
			PageSize: req.PageSize,
		},
		Username: req.Username,
		Gender:   req.Gender,
		Enable:   req.Enable,
	})
	if errCode != 0 {
		handler.ErrorLog(
			logger.AdminLogger,
			"adminSvc.ListPage 调用失败",
			errCode,
			err,
			zap.Uint64("adminID", adminID),
			zap.Int("req.page", req.Page),
			zap.Int("req.pageSize", req.PageSize),
			zap.Any("req.username", req.Username),
			zap.Any("req.gender", req.Gender),
			zap.Any("req.enable", req.Enable),
		)
		response.Error(c, lang, errCode)
		return
	}
	// 返回结果
	response.Success(c, lang, resp.AdminListPageResp{
		Total:    svcResp.Total,
		PageData: toAdminListItems(svcResp.PageData),
	})
}

// Details 处理获取管理员详情请求
// POST /api/admin/admin/detail
func (h *Handler) Details(c *gin.Context) {
	// 获取上下文/语言配置
	ctx := c.Request.Context()
	lang := i18n.GetLang(ctx)
	// 获取管理员ID
	adminID, ok := handler.GetAdminID(c)
	if !ok {
		response.Error(c, lang, 20101)
		return
	}
	// 执行请求
	svcResp, errCode, err := h.adminSvc.Details(ctx, adminID)
	if errCode != 0 {
		handler.ErrorLog(
			logger.AdminLogger,
			"adminSvc.Details 调用失败",
			errCode,
			err,
			zap.Uint64("adminID", adminID),
		)
		response.Error(c, lang, errCode)
		return
	}
	// 返回结果
	response.Success(c, lang, resp.AdminDetailsResp{
		ID:        svcResp.ID,
		Username:  svcResp.Username,
		Enable:    svcResp.Enable,
		CreatedAt: svcResp.CreatedAt,
		UpdatedAt: svcResp.UpdatedAt,
		Profile: resp.AdminDetailsProfileItem{
			ID:       svcResp.Profile.ID,
			Nickname: svcResp.Profile.Nickname,
			Gender:   svcResp.Profile.Gender,
			Avatar:   svcResp.Profile.Avatar,
			Address:  svcResp.Profile.Address,
			Email:    svcResp.Profile.Email,
		},
		Roles: toAdminDetailsRoleItem(svcResp.Roles),
		CurrentRole: resp.AdminDetailsRoleItem{
			ID:     svcResp.CurrentRole.ID,
			Code:   svcResp.CurrentRole.Code,
			Name:   svcResp.CurrentRole.Name,
			Enable: svcResp.CurrentRole.Enable,
		},
	})
}

// Save 处理创建或变更管理员账号请求
// POST /api/admin/admin/save
func (h *Handler) Save(c *gin.Context) {
	// 获取上下文/语言配置
	ctx := c.Request.Context()
	lang := i18n.GetLang(ctx)
	// 获取管理员ID
	adminID, ok := handler.GetAdminID(c)
	if !ok {
		response.Error(c, lang, 20101)
		return
	}
	// 获取请求参数
	var req input.AdminSaveReq
	if code, ok, err := handler.BindAndValidate(c, &req); !ok {
		handler.ErrorLog(
			logger.AdminLogger,
			"Save 参数异常",
			code,
			err,
		)
		response.Error(c, lang, code)
		return
	}
	// 执行请求
	errCode, err := h.adminSvc.Save(ctx, admin.SaveReq{
		ID:       req.ID,
		Enable:   req.Enable,
		Username: req.Username,
		Password: req.Password,
		RoleIds:  req.RoleIds,
	})
	if errCode != 0 {
		handler.ErrorLog(
			logger.AdminLogger,
			"adminSvc.Save 调用失败",
			errCode,
			err,
			zap.Uint64("adminID", adminID),
			zap.Any("req.id", req.ID),
			zap.String("req.username", req.Username),
			zap.Any("req.role_ids", req.RoleIds),
		)
		response.Error(c, lang, errCode)
		return
	}
	// 返回结果
	response.Success(c, lang, nil)
}

// Delete 处理删除管理员请求
// POST /api/admin/admin/delete
func (h *Handler) Delete(c *gin.Context) {
	// 获取上下文/语言配置
	ctx := c.Request.Context()
	lang := i18n.GetLang(ctx)
	// 获取管理员ID
	adminID, ok := handler.GetAdminID(c)
	if !ok {
		response.Error(c, lang, 20101)
		return
	}
	// 获取请求参数
	var req input.AdminDeleteReq
	if code, ok, err := handler.BindAndValidate(c, &req); !ok {
		handler.ErrorLog(
			logger.AdminLogger,
			"Delete 参数异常",
			code,
			err,
		)
		response.Error(c, lang, code)
		return
	}
	// 执行请求
	errCode, err := h.adminSvc.Delete(ctx, req.ID)
	if errCode != 0 {
		handler.ErrorLog(
			logger.AdminLogger,
			"adminSvc.Delete 调用失败",
			errCode,
			err,
			zap.Uint64("adminID", adminID),
		)
		response.Error(c, lang, errCode)
		return
	}
	// 返回结果
	response.Success(c, lang, nil)
}

// UpdatePassword 处理变更管理员账号密码请求
// POST /api/admin/admin/update-password
func (h *Handler) UpdatePassword(c *gin.Context) {
	// 获取上下文/语言配置
	ctx := c.Request.Context()
	lang := i18n.GetLang(ctx)
	// 获取管理员ID
	adminID, ok := handler.GetAdminID(c)
	if !ok {
		response.Error(c, lang, 20101)
		return
	}
	// 获取请求参数
	var req input.AdminResetAdminPasswordReq
	if code, ok, err := handler.BindAndValidate(c, &req); !ok {
		handler.ErrorLog(
			logger.AdminLogger,
			"UpdatePassword 参数异常",
			code,
			err,
		)
		response.Error(c, lang, code)
		return
	}
	// 执行请求
	errCode, err := h.adminSvc.ResetAdminPassword(ctx, req.ID, req.Password)
	if errCode != 0 {
		handler.ErrorLog(
			logger.AdminLogger,
			"adminSvc.ResetAdminPassword 调用失败",
			errCode,
			err,
			zap.Uint64("adminID", adminID),
			zap.Uint64("req.id", req.ID),
		)
		response.Error(c, lang, errCode)
		return
	}
	// 返回结果
	response.Success(c, lang, nil)
}

// UpdatePassword 处理变更管理员个人信息请求
// POST /api/admin/admin/update-profile
func (h *Handler) UpdateProfile(c *gin.Context) {
	// 获取上下文/语言配置
	ctx := c.Request.Context()
	lang := i18n.GetLang(ctx)
	// 获取管理员ID
	adminID, ok := handler.GetAdminID(c)
	if !ok {
		response.Error(c, lang, 20101)
		return
	}
	// 获取请求参数
	var req input.AdminUpdateProfileReq
	if code, ok, err := handler.BindAndValidate(c, &req); !ok {
		handler.ErrorLog(
			logger.AdminLogger,
			"UpdateProfile 参数异常",
			code,
			err,
		)
		response.Error(c, lang, code)
		return
	}
	// 执行请求
	errCode, err := h.adminSvc.UpdateProfile(ctx, admin.UpdateProfileReq{
		ID:       req.ID,
		Nickname: req.Nickname,
		Gender:   req.Gender,
		Address:  req.Address,
		Email:    req.Email,
	})
	if errCode != 0 {
		handler.ErrorLog(
			logger.AdminLogger,
			"adminSvc.UpdateProfile 调用失败",
			errCode,
			err,
			zap.Uint64("adminID", adminID),
			zap.Uint64("req.id", req.ID),
			zap.String("req.nickname", req.Nickname),
			zap.Int("req.gender", req.Gender),
			zap.Any("req.Address", req.Address),
			zap.Any("req.Email", req.Email),
		)
		response.Error(c, lang, errCode)
		return
	}
	// 返回结果
	response.Success(c, lang, nil)
}
