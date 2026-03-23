package admin

import (
	"context"

	"github.com/redis/go-redis/v9"
	"github.com/zxc7563598/oneadmin/internal/enum"
	"github.com/zxc7563598/oneadmin/internal/model"
	"github.com/zxc7563598/oneadmin/internal/repository/admin"
	"github.com/zxc7563598/oneadmin/internal/repository/admin_role"
	"github.com/zxc7563598/oneadmin/internal/repository/role"
	"github.com/zxc7563598/oneadmin/pkg/crypto"
	"github.com/zxc7563598/oneadmin/pkg/jwt"
	"github.com/zxc7563598/oneadmin/pkg/ptr"
	"github.com/zxc7563598/oneadmin/pkg/timeutil"
	"gorm.io/gorm"
)

type Service struct {
	db            *gorm.DB
	adminRepo     admin.Repository
	adminRoleRepo admin_role.Repository
	roleRepo      role.Repository
	rdb           *redis.Client
}

func New(adminRepo admin.Repository, adminRoleRepo admin_role.Repository, roleRepo role.Repository, db *gorm.DB, rdb *redis.Client) *Service {
	return &Service{
		adminRepo:     adminRepo,
		adminRoleRepo: adminRoleRepo,
		roleRepo:      roleRepo,
		db:            db,
		rdb:           rdb,
	}
}

// Login 用于处理管理员登录
func (s *Service) Login(ctx context.Context, username, password string, captcha *string) (LoginResp, int, error) {
	// 获取管理员信息
	admin, err := s.adminRepo.GetByUsername(ctx, nil, username)
	if err != nil {
		return LoginResp{}, 60101, err
	}
	// 验证密码是否正确
	if admin == nil {
		return LoginResp{}, 50101, nil
	}
	if !crypto.CheckPassword(admin.Password, password) {
		return LoginResp{}, 40101, nil
	}
	// 验证账号是否启用
	if admin.Enable != enum.EnableEnable {
		return LoginResp{}, 40102, nil
	}
	// 获取角色code
	role, err := s.roleRepo.GetByID(ctx, nil, admin.RoleID)
	if err != nil {
		return LoginResp{}, 60101, err
	}
	if role == nil {
		return LoginResp{}, 50102, nil
	}
	// 更新token
	updateToken, errCode, err := s.updateToken(ctx, admin.ID, admin.RoleID, role.Code)
	if err != nil {
		return LoginResp{}, errCode, err
	}
	// 返回数据
	return LoginResp{
		AccessToken:  updateToken.AccessToken,
		RefreshToken: updateToken.RefreshToken,
	}, 0, nil
}

// RefreshLogin 用于刷新管理员登录状态
func (s *Service) RefreshLogin(ctx context.Context, refreshToken string) (RefreshLoginResp, int, error) {
	claims, err := jwt.ParseToken(refreshToken)
	if err != nil {
		return RefreshLoginResp{}, 10002, err
	}
	if claims.Type != "refresh" {
		return RefreshLoginResp{}, 10003, nil
	}
	// 获取管理员信息
	admin, err := s.adminRepo.GetByID(ctx, nil, claims.ID)
	if err != nil {
		return RefreshLoginResp{}, 60101, err
	}
	// 验证信息
	if admin == nil {
		return RefreshLoginResp{}, 50101, nil
	}
	if admin.Token == nil || *admin.Token != refreshToken {
		return RefreshLoginResp{}, 20001, nil
	}
	// 获取角色code
	role, err := s.roleRepo.GetByID(ctx, nil, admin.RoleID)
	if err != nil {
		return RefreshLoginResp{}, 60101, err
	}
	if role == nil {
		return RefreshLoginResp{}, 50102, nil
	}
	// 更新token
	updateToken, errCode, err := s.updateToken(ctx, claims.ID, admin.RoleID, role.Code)
	if err != nil {
		return RefreshLoginResp{}, errCode, err
	}
	// 返回数据
	return RefreshLoginResp{
		AccessToken:  updateToken.AccessToken,
		RefreshToken: updateToken.RefreshToken,
	}, 0, nil
}

// Logout 用于退出管理员登录
func (s *Service) Logout(ctx context.Context, adminID uint64) (int, error) {
	// 清空用户token
	if s.rdb != nil {
		err := s.rdb.Del(ctx,
			jwt.AdminTokenKey(adminID),
			jwt.AdminRefreshKey(adminID),
		).Err()
		if err != nil {
			return 60107, err
		}
	}
	if err := s.adminRepo.UpdateToken(ctx, nil, adminID, nil); err != nil {
		return 60104, err
	}
	// 返回数据
	return 0, nil
}

// SwitchRole 用于切换管理员角色信息
func (s *Service) SwitchRole(ctx context.Context, adminID, roleID uint64) (SwitchRoleResp, int, error) {
	// 获取用户角色是否存在
	exists, err := s.adminRoleRepo.HasRole(ctx, nil, adminID, roleID)
	if err != nil {
		return SwitchRoleResp{}, 60101, err
	}
	if !exists {
		return SwitchRoleResp{}, 30101, nil
	}
	// 获取角色信息
	role, err := s.roleRepo.GetByID(ctx, nil, roleID)
	if err != nil {
		return SwitchRoleResp{}, 60101, err
	}
	if role == nil {
		return SwitchRoleResp{}, 50102, nil
	}
	if role.Enable != enum.EnableEnable {
		return SwitchRoleResp{}, 40103, nil
	}
	// 变更角色
	if err := s.adminRepo.UpdateRole(ctx, nil, adminID, roleID); err != nil {
		return SwitchRoleResp{}, 60108, err
	}
	// 更新token
	updateToken, errCode, err := s.updateToken(ctx, adminID, roleID, role.Code)
	if err != nil {
		return SwitchRoleResp{}, errCode, err
	}
	// 返回数据
	return SwitchRoleResp{
		AccessToken:  updateToken.AccessToken,
		RefreshToken: updateToken.RefreshToken,
	}, 0, nil
}

// ChangePassword 用于根据管理员旧密码修改密码
func (s *Service) ChangePassword(ctx context.Context, adminID uint64, oldPassword, newPassword string) (int, error) {
	// 获取管理员信息
	admin, err := s.adminRepo.GetByID(ctx, nil, adminID)
	if err != nil {
		return 60101, err
	}
	if admin == nil {
		return 50101, nil
	}
	// 验证旧密码
	if !crypto.CheckPassword(admin.Password, oldPassword) {
		return 40101, nil
	}
	// 防止新旧密码相同
	if crypto.CheckPassword(admin.Password, newPassword) {
		return 40104, nil
	}
	password, err := crypto.HashPassword(newPassword)
	if err != nil {
		return 60109, err
	}
	if err := s.adminRepo.UpdatePassword(ctx, nil, adminID, password); err != nil {
		return 60110, err
	}
	// 清空管理员登录状态并返回
	return s.Logout(ctx, adminID)
}

// ListPage 用于获取管理员列表信息
func (s *Service) ListPage(ctx context.Context, req ListPageReq) (ListPageResp, int, error) {
	// 获取列表数据
	offset, limit := req.OffsetLimit()
	admins, total, err := s.adminRepo.ListPage(ctx, nil, model.AdminListQuery{
		Username: req.Username,
		Gender:   req.Gender,
		Enable:   req.Enable,
		Offset:   offset,
		Limit:    limit,
	})
	if err != nil {
		return ListPageResp{}, 60101, err
	}
	// 获取列表管理员角色
	adminIDs := make([]uint64, 0, len(admins))
	for _, v := range admins {
		adminIDs = append(adminIDs, v.ID)
	}
	adminRoles, err := s.adminRoleRepo.GetByRoleIDs(ctx, nil, adminIDs)
	if err != nil {
		return ListPageResp{}, 60101, err
	}
	// 获取全部角色
	roles, err := s.roleRepo.FindAll(ctx, nil)
	if err != nil {
		return ListPageResp{}, 60101, err
	}
	// 组装数据
	roleMap := make(map[uint64]DetailsRoleItem)
	for _, v := range roles {
		roleMap[v.ID] = DetailsRoleItem{
			ID:     v.ID,
			Code:   v.Code,
			Name:   v.Name,
			Enable: v.Enable == enum.EnableEnable,
		}
	}
	//
	adminRoleMap := make(map[uint64][]DetailsRoleItem)
	for _, ar := range adminRoles {
		role, ok := roleMap[ar.RoleID]
		if !ok {
			continue
		}
		adminRoleMap[ar.AdminID] = append(adminRoleMap[ar.AdminID], role)
	}
	respList := make([]ListPageItem, 0, len(admins))
	for _, v := range admins {
		item := ListPageItem{
			ID:        v.ID,
			Username:  v.Username,
			Enable:    v.Enable == enum.EnableEnable,
			Gender:    int(v.Gender),
			Avatar:    v.Avatar,
			Address:   v.Address,
			Email:     v.Email,
			Roles:     adminRoleMap[v.ID],
			CreatedAt: timeutil.Format(v.CreatedAt),
			UpdatedAt: timeutil.Format(v.UpdatedAt),
		}
		respList = append(respList, item)
	}
	// 返回数据
	return ListPageResp{
		Total:    total,
		PageData: respList,
	}, 0, nil
}

// Details 用于获取单个管理员的详细信息
func (s *Service) Details(ctx context.Context, adminID uint64) (DetailsResp, int, error) {
	// 获取管理员信息
	admin, err := s.adminRepo.GetByID(ctx, nil, adminID)
	if err != nil {
		return DetailsResp{}, 60101, err
	}
	if admin == nil {
		return DetailsResp{}, 50101, nil
	}
	// 获取角色信息
	roles, err := s.roleRepo.FindEnabled(ctx, nil)
	if err != nil {
		return DetailsResp{}, 60101, err
	}
	// 组装参数
	roleList := make([]DetailsRoleItem, 0, len(roles))
	var currentRole DetailsRoleItem
	for _, v := range roles {
		item := DetailsRoleItem{
			ID:     v.ID,
			Code:   v.Code,
			Name:   v.Name,
			Enable: v.Enable == enum.EnableEnable,
		}
		roleList = append(roleList, item)
		if v.ID == admin.RoleID {
			currentRole = item
		}
	}
	return DetailsResp{
		ID:        admin.ID,
		Username:  admin.Username,
		Enable:    admin.Enable == enum.EnableEnable,
		CreatedAt: timeutil.Format(admin.CreatedAt),
		UpdatedAt: timeutil.Format(admin.UpdatedAt),
		Profile: DetailsProfileItem{
			ID:       admin.ID,
			Nickname: admin.Nickname,
			Gender:   int(admin.Gender),
			Avatar:   admin.Avatar,
			Address:  ptr.Deref(admin.Address),
			Email:    ptr.Deref(admin.Email),
		},
		Roles:       roleList,
		CurrentRole: currentRole,
	}, 0, nil
}

// Save 用于创建或修改管理员信息
func (s *Service) Save(ctx context.Context, req SaveReq) (int, error) {
	// 开启事务
	var adminID uint64
	var errCode int
	var err error
	if req.ID == nil {
		adminID, errCode, err = s.add(ctx, req)
		if errCode > 0 {
			return errCode, err
		}
	} else {
		adminID, errCode, err = s.update(ctx, req)
		if errCode > 0 {
			return errCode, err
		}
	}
	errCode, err = s.bindRole(ctx, adminID, req.RoleIds)
	if errCode > 0 {
		return errCode, err
	}
	return 0, nil
}

func (s *Service) add(ctx context.Context, req SaveReq) (uint64, int, error) {
	if req.Enable == nil {
		return 0, 10101, nil
	}
	enable := enum.EnableDisable
	if *req.Enable {
		enable = enum.EnableEnable
	}
	if req.Password == nil {
		return 0, 10103, nil
	}
	if len(*req.Password) < 6 {
		return 0, 10104, nil
	}
	if len(*req.Password) > 32 {
		return 0, 10105, nil
	}
	if len(req.RoleIds) > 0 {
		return 0, 10108, nil
	}
	// 添加数据
	password, err := crypto.HashPassword(*req.Password)
	if err != nil {
		return 0, 60109, err
	}
	admin, err := s.adminRepo.Create(ctx, nil, &model.Admin{
		Nickname: req.Username,
		Username: req.Username,
		Password: password,
		RoleID:   req.RoleIds[0],
		Gender:   enum.GenderUnknown,
		Enable:   enable,
	})
	if err != nil {
		return 0, 60111, err
	}
	return admin.ID, 0, nil
}

func (s *Service) update(ctx context.Context, req SaveReq) (uint64, int, error) {
	var enable *enum.Enable
	if req.Enable != nil {
		if *req.Enable {
			val := enum.EnableEnable
			enable = &val
		} else {
			val := enum.EnableDisable
			enable = &val
		}
	}
	var roleID *uint64
	if len(req.RoleIds) > 0 {
		roleID = &req.RoleIds[0]
	}
	// 变更数据
	if err := s.adminRepo.UpdateInfo(ctx, nil, *req.ID, model.AdminUpdateInfoForm{
		Username: req.Username,
		Enable:   enable,
		RoleID:   roleID,
	}); err != nil {
		return 0, 60111, err
	}
	return *req.ID, 0, nil
}

func (s *Service) bindRole(ctx context.Context, adminID uint64, roleIds []uint64) (int, error) {
	if len(roleIds) > 0 {
		// 删除原有角色信息
		if err := s.adminRoleRepo.DeleteByAdminID(ctx, nil, adminID); err != nil {
			return 60114, err
		}
		// 重新绑定角色信息
		adminRoleList := make([]model.AdminRole, 0, len(roleIds))
		for _, v := range roleIds {
			adminRoleList = append(adminRoleList, model.AdminRole{
				AdminID: adminID,
				RoleID:  v,
			})
		}
		if err := s.adminRoleRepo.CreateBatch(ctx, nil, adminRoleList); err != nil {
			return 60115, err
		}
	}
	return 0, nil
}

// Delete 用于删除管理员信息
func (s *Service) Delete(ctx context.Context, adminID uint64) (int, error) {
	// 获取管理员信息
	admin, err := s.adminRepo.GetByID(ctx, nil, adminID)
	if err != nil {
		return 60101, err
	}
	if admin == nil {
		return 50101, nil
	}
	// 开启事务，删除管理员
	err = s.db.Transaction(func(tx *gorm.DB) error {
		if err := s.adminRoleRepo.DeleteByAdminID(ctx, tx, adminID); err != nil {
			return err
		}
		if err := s.adminRepo.Delete(ctx, tx, adminID); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return 60112, err
	}
	// 返回数据
	return 0, nil
}

// ResetAdminPassword 修改管理员密码
func (s *Service) ResetAdminPassword(ctx context.Context, adminID uint64, newPassword string) (int, error) {
	// 验证管理员是否存在
	admin, err := s.adminRepo.GetByID(ctx, nil, adminID)
	if err != nil {
		return 60101, err
	}
	if admin == nil {
		return 50101, nil
	}
	password, err := crypto.HashPassword(newPassword)
	if err != nil {
		return 60109, err
	}
	if err := s.adminRepo.UpdatePassword(ctx, nil, adminID, password); err != nil {
		return 60110, err
	}
	// 清空管理员登录状态并返回
	return s.Logout(ctx, adminID)
}

// UpdateProfile 修改管理员个人信息
func (s *Service) UpdateProfile(ctx context.Context, req UpdateProfileReq) (int, error) {
	// 验证管理员是否存在
	admin, err := s.adminRepo.GetByID(ctx, nil, req.ID)
	if err != nil {
		return 60101, err
	}
	if admin == nil {
		return 50101, nil
	}
	// 变更管理员信息
	if err := s.adminRepo.UpdateProfile(ctx, nil, req.ID, model.AdminUpdateProfileForm{
		Nickname: req.Nickname,
		Email:    req.Email,
		Address:  req.Address,
		Gender:   req.Gender,
		Avatar:   req.Avatar,
	}); err != nil {
		return 60113, err
	}
	return 0, nil
}

// updateToken 更新token
func (s *Service) updateToken(ctx context.Context, adminID, roleID uint64, roleCode string) (TokenResp, int, error) {
	accessToken, err := jwt.GenerateAccessToken(adminID, "admin", roleID, roleCode)
	if err != nil {
		return TokenResp{}, 60102, err
	}
	newRefreshToken, err := jwt.GenerateRefreshToken(adminID, "admin", roleID, roleCode)
	if err != nil {
		return TokenResp{}, 60103, err
	}
	if err := s.adminRepo.UpdateToken(ctx, nil, adminID, &newRefreshToken); err != nil {
		return TokenResp{}, 60104, err
	}
	if s.rdb != nil {
		if err := s.rdb.Set(ctx, jwt.AdminTokenKey(adminID), accessToken, jwt.AccessTTL()).Err(); err != nil {
			return TokenResp{}, 60105, err
		}
		if err := s.rdb.Set(ctx, jwt.AdminRefreshKey(adminID), newRefreshToken, jwt.RefreshTTL()).Err(); err != nil {
			return TokenResp{}, 60106, err
		}
	}
	return TokenResp{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
	}, 0, nil
}
