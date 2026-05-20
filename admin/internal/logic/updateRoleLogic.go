// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"
	"errors"
	"regexp"

	"myproject/admin/internal/svc"
	"myproject/admin/internal/types"
	"myproject/admin/internal/util"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateRoleLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewUpdateRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateRoleLogic {
	return &UpdateRoleLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *UpdateRoleLogic) UpdateRole(req *types.UpdateRoleRequest) (resp *types.RoleInfo, err error) {
	// 参数校验
	// 角色名称：不超过20个字符
	if req.Name != "" && len(req.Name) > 20 {
		return nil, errors.New("角色名称不超过20个字符")
	}

	// 角色编码：支持大小写字母、数字和下划线，不超过20个字符
	if req.Code != "" {
		codePattern := regexp.MustCompile(`^[A-Za-z0-9_]{1,20}$`)
		if !codePattern.MatchString(req.Code) {
			return nil, errors.New("角色编码支持大小写字母、数字和下划线，不超过20个字符")
		}
	}

	// 角色描述：不超过200个字符
	if req.Desc != "" && len(req.Desc) > 200 {
		return nil, errors.New("角色描述不超过200个字符")
	}

	// 1. 根据ID获取角色
	role, err := l.svcCtx.RoleModel.GetByID(req.ID)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, errors.New("角色不存在")
	}

	// 2. 更新角色信息
	if req.Name != "" {
		role.Name = req.Name
	}

	if req.Code != "" {
		// 检查角色代码是否已存在
		existingRole, err := l.svcCtx.RoleModel.GetByCode(req.Code)
		if err != nil {
			return nil, err
		}
		if existingRole != nil && existingRole.ID != req.ID {
			return nil, errors.New("角色代码已存在")
		}
		role.Code = req.Code
	}

	if req.Desc != "" {
		role.Desc = req.Desc
	}

	// 状态字段：0 表示禁用，1 表示启用，需要始终更新
	role.Status = req.Status

	// 3. 保存更新
	err = l.svcCtx.RoleModel.Update(role)
	if err != nil {
		return nil, err
	}

	// 4. 构建响应
	resp = &types.RoleInfo{
		ID:          role.ID,
		Name:        role.Name,
		Code:        role.Code,
		Desc:        role.Desc,
		Status:      role.Status,
		Permissions: util.ConvertPermissions(role.Permissions),
	}

	return resp, nil
}
