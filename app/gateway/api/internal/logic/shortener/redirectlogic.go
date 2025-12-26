// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package shortener

import (
	"context"
	"errors"
	"go-link/app/shortener/rpc/shortener"

	"go-link/app/gateway/api/internal/svc"
	"go-link/app/gateway/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type RedirectLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRedirectLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RedirectLogic {
	return &RedirectLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RedirectLogic) Redirect(req *types.GetLinkReq) (resp *types.GetLinkResp, err error) {
	// 调用rpc获取原始路径
	rpcResp, err := l.svcCtx.ShortenerRpc.GetOrigin(l.ctx, &shortener.GetOriginReq{
		ShortCode: req.Code,
	})

	if err != nil {
		l.Error("跳转错误:", err)
		return nil, errors.New("跳转错误")
	}

	return &types.GetLinkResp{
		Url: rpcResp.OriginUrl,
	}, nil
}
