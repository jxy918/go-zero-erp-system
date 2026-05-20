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

type CreateUserLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewCreateUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateUserLogic {
	return &CreateUserLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *CreateUserLogic) CreateUser(req *types.CreateUserRequest) (resp *types.UserInfo, err error) {
	// 参数校验
	// 用户名：只支持英文大小写和数字，不超过20个字符
	usernamePattern := regexp.MustCompile(`^[A-Za-z0-9]{1,20}$`)
	if !usernamePattern.MatchString(req.Username) {
		return nil, errors.New("用户名只支持英文大小写和数字，不超过20个字符")
	}

	// 昵称：不超过20个字符
	if len(req.Nickname) > 20 {
		return nil, errors.New("昵称不超过20个字符")
	}

	// 邮箱格式校验
	emailPattern := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if req.Email != "" && !emailPattern.MatchString(req.Email) {
		return nil, errors.New("邮箱格式不正确")
	}

	// 手机格式校验
	phonePattern := regexp.MustCompile(`^1[3-9]\d{9}$`)
	if req.Phone != "" && !phonePattern.MatchString(req.Phone) {
		return nil, errors.New("手机号码格式不正确")
	}

	// 1. 检查用户名是否已存在
	existingUser, err := l.svcCtx.UserModel.GetByUsername(req.Username)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, errors.New("用户名已存在")
	}

	// 2. 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// 3. 创建用户
	user := &model.User{
		Username: req.Username,
		Password: string(hashedPassword),
		Nickname: req.Nickname,
		Email:    req.Email,
		Phone:    req.Phone,
		Status:   req.Status,
	}

	err = l.svcCtx.UserModel.Create(user)
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
		Roles:    []types.RoleInfo{},
	}

	return resp, nil
}
