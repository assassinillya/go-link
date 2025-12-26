// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"
	"go-link/app/user/rpc/user"

	"go-link/app/gateway/api/internal/svc"
	"go-link/app/gateway/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(req *types.LoginReq) (resp *types.LoginResp, err error) {
	// 打印个日志
	logx.Infof("正在请求登录，用户名：%s", req.Username)

	// 调用rpc
	rpcResp, err := l.svcCtx.UserRpc.Login(l.ctx, &user.LoginReq{
		Username: req.Username,
		Password: req.Password,
	})

	if err != nil {
		return nil, err
	}

	return &types.LoginResp{
		Token: rpcResp.Token,
		Msg:   rpcResp.Msg,
	}, nil
}
