// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"
	"errors"

	"myproject/admin/internal/metric"
	"myproject/admin/internal/model"
	"myproject/admin/internal/svc"
	"myproject/admin/internal/types"
	"myproject/admin/internal/util"

	"github.com/zeromicro/go-zero/core/logx"
	"golang.org/x/crypto/bcrypt"
)

type LoginLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(req *types.LoginRequest) (resp *types.LoginResponse, err error) {
	user, err := l.svcCtx.UserModel.GetByUsername(req.Username)
	if err != nil {
		if metric.IsEnabled() {
			metric.LoginCounter.Inc("failure")
		}
		return nil, err
	}
	if user == nil {
		if metric.IsEnabled() {
			metric.LoginCounter.Inc("failure")
		}
		return nil, errors.New("用户不存在")
	}

	if user.Status == 0 {
		if metric.IsEnabled() {
			metric.LoginCounter.Inc("failure")
		}
		return nil, errors.New("用户已被禁用")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		if metric.IsEnabled() {
			metric.LoginCounter.Inc("failure")
		}
		return nil, errors.New("密码错误")
	}

	// 登录成功
	if metric.IsEnabled() {
		metric.LoginCounter.Inc("success")
	}

	var roleID uint
	var isAdmin bool
	if len(user.Roles) > 0 {
		roleID = user.Roles[0].ID
		isAdmin = user.Roles[0].Code == "admin"
	}

	token, err := util.GenerateToken(user.ID, user.Username, roleID, isAdmin, l.svcCtx.Config.JWT.AccessSecret, l.svcCtx.Config.JWT.AccessExpire)
	if err != nil {
		return nil, err
	}

	resp = &types.LoginResponse{
		Token: token,
		User: types.UserInfo{
			ID:       user.ID,
			Username: user.Username,
			Nickname: user.Nickname,
			Email:    user.Email,
			Phone:    user.Phone,
			Status:   user.Status,
			Roles:    l.convertRoles(user.Roles),
		},
	}

	return resp, nil
}

func (l *LoginLogic) convertRoles(roles []model.Role) []types.RoleInfo {
	var roleInfos []types.RoleInfo
	for _, role := range roles {
		roleInfos = append(roleInfos, types.RoleInfo{
			ID:          role.ID,
			Name:        role.Name,
			Code:        role.Code,
			Desc:        role.Desc,
			Status:      role.Status,
			Permissions: l.convertPermissions(role.Permissions),
			Menus:       l.convertMenus(role.Menus),
		})
	}
	return roleInfos
}

func (l *LoginLogic) convertMenus(menus []model.Menu) []types.MenuInfo {
	var menuInfos []types.MenuInfo
	for _, menu := range menus {
		menuInfos = append(menuInfos, types.MenuInfo{
			ID:        menu.ID,
			Name:      menu.Name,
			Code:      menu.Code,
			Desc:      menu.Desc,
			ParentID:  menu.ParentID,
			Path:      menu.Path,
			Component: menu.Component,
			Icon:      menu.Icon,
			Sort:      menu.Sort,
			Status:    menu.Status,
		})
	}
	return menuInfos
}

func (l *LoginLogic) convertPermissions(permissions []model.Permission) []types.PermissionInfo {
	var permissionInfos []types.PermissionInfo
	for _, permission := range permissions {
		permissionInfos = append(permissionInfos, types.PermissionInfo{
			ID:     permission.ID,
			Name:   permission.Name,
			Code:   permission.Code,
			Desc:   permission.Desc,
			Sort:   permission.Sort,
			Status: permission.Status,
		})
	}
	return permissionInfos
}
