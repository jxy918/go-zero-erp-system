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
)

type CreateRoleLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewCreateRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateRoleLogic {
	return &CreateRoleLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *CreateRoleLogic) CreateRole(req *types.CreateRoleRequest) (resp *types.RoleInfo, err error) {
	// 参数校验
	// 角色名称：不超过20个字符
	if len(req.Name) > 20 {
		return nil, errors.New("角色名称不超过20个字符")
	}

	// 角色编码：支持大小写字母、数字和下划线，不超过20个字符
	codePattern := regexp.MustCompile(`^[A-Za-z0-9_]{1,20}$`)
	if !codePattern.MatchString(req.Code) {
		return nil, errors.New("角色编码支持大小写字母、数字和下划线，不超过20个字符")
	}

	// 角色描述：不超过200个字符
	if len(req.Desc) > 200 {
		return nil, errors.New("角色描述不超过200个字符")
	}

	// 1. 检查角色代码是否已存在
	existingRole, err := l.svcCtx.RoleModel.GetByCode(req.Code)
	if err != nil {
		return nil, err
	}
	if existingRole != nil {
		return nil, errors.New("角色代码已存在")
	}

	// 2. 创建角色
	role := &model.Role{
		Name:   req.Name,
		Code:   req.Code,
		Desc:   req.Desc,
		Status: req.Status,
	}

	err = l.svcCtx.RoleModel.Create(role)
	if err != nil {
		return nil, err
	}

	// 3. 构建响应
	resp = &types.RoleInfo{
		ID:          role.ID,
		Name:        role.Name,
		Code:        role.Code,
		Desc:        role.Desc,
		Status:      role.Status,
		Permissions: []types.PermissionInfo{},
	}

	return resp, nil
}
