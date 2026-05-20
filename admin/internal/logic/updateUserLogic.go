// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"
	"errors"
	"regexp"

	"myproject/admin/internal/model"
	"myproject/admin/internal/svc"
	"myproject/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	"golang.org/x/crypto/bcrypt"
)

type UpdateUserLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewUpdateUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateUserLogic {
	return &UpdateUserLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *UpdateUserLogic) UpdateUser(req *types.UpdateUserRequest) (resp *types.UserInfo, err error) {
	// 参数校验
	// 用户名：只支持英文大小写和数字，不超过20个字符
	if req.Username != "" {
		usernamePattern := regexp.MustCompile(`^[A-Za-z0-9]{1,20}$`)
		if !usernamePattern.MatchString(req.Username) {
			return nil, errors.New("用户名只支持英文大小写和数字，不超过20个字符")
		}
	}

	// 昵称：不超过20个字符
	if req.Nickname != "" && len(req.Nickname) > 20 {
		return nil, errors.New("昵称不超过20个字符")
	}

	// 邮箱格式校验
	if req.Email != "" {
		emailPattern := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
		if !emailPattern.MatchString(req.Email) {
			return nil, errors.New("邮箱格式不正确")
		}
	}

	// 手机格式校验
	if req.Phone != "" {
		phonePattern := regexp.MustCompile(`^1[3-9]\d{9}$`)
		if !phonePattern.MatchString(req.Phone) {
			return nil, errors.New("手机号码格式不正确")
		}
	}

	// 1. 根据ID获取用户
	user, err := l.svcCtx.UserModel.GetByID(req.ID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("用户不存在")
	}

	// 2. 更新用户信息
	if req.Username != "" {
		// 检查用户名是否已存在
		existingUser, err := l.svcCtx.UserModel.GetByUsername(req.Username)
		if err != nil {
			return nil, err
		}
		if existingUser != nil && existingUser.ID != req.ID {
			return nil, errors.New("用户名已存在")
		}
		user.Username = req.Username
	}

	if req.Password != "" {
		// 加密密码
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		user.Password = string(hashedPassword)
	}

	if req.Nickname != "" {
		user.Nickname = req.Nickname
	}

	if req.Email != "" {
		user.Email = req.Email
	}

	if req.Phone != "" {
		user.Phone = req.Phone
	}

	// 状态字段：0 表示禁用，1 表示启用，需要始终更新
	user.Status = req.Status

	// 3. 保存更新
	err = l.svcCtx.UserModel.Update(user)
	if err != nil {
		return nil, err
	}

	// 4. 构建响应
	resp = &types.UserInfo{
		ID:       user.ID,
		Username: user.Username,
		Nickname: user.Nickname,
		Email:    user.Email,
		Phone:    user.Phone,
		Status:   user.Status,
		Roles:    l.convertRoles(user.Roles),
	}

	return resp, nil
}

// convertRoles 转换角色列表
func (l *UpdateUserLogic) convertRoles(roles []model.Role) []types.RoleInfo {
	var roleInfos []types.RoleInfo
	for _, role := range roles {
		roleInfos = append(roleInfos, types.RoleInfo{
			ID:          role.ID,
			Name:        role.Name,
			Code:        role.Code,
			Desc:        role.Desc,
			Status:      role.Status,
			Permissions: l.convertPermissions(role.Permissions),
		})
	}
	return roleInfos
}

// convertPermissions 转换权限列表
func (l *UpdateUserLogic) convertPermissions(permissions []model.Permission) []types.PermissionInfo {
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
