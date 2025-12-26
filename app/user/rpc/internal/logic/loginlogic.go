package logic

import (
	"context"
	"errors"
	"go-link/app/user/model"
	"go-link/app/user/rpc/internal/svc"
	"go-link/app/user/rpc/user"
	"go-link/pkg/utils"
	"gorm.io/gorm"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *LoginLogic) Login(in *user.LoginReq) (*user.LoginResp, error) {

	username := in.Username
	password := in.Password

	var u model.User
	err := l.svcCtx.DB.Where("username = ?", username).Take(&u).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户不存在") // rpc直接返回error就行了
			// 不需要返回http状态码
		}
		return nil, errors.New("数据库查询失败")
	}

	// 校验密码
	if !utils.CheckPassword(password, u.Password) {
		return nil, errors.New("密码错误")
	}

	// 生成Token
	token, err := utils.GenerateToken(u.ID)
	if err != nil {
		return nil, errors.New("Token 生成失败")
	}

	return &user.LoginResp{
		Token: token,
		Msg:   "登录成功",
	}, nil

}
