// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"
	"encoding/json"
	"errors"
	"go-link/app/shortener/rpc/shortener"

	"go-link/app/gateway/api/internal/svc"
	"go-link/app/gateway/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateLinkLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateLinkLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateLinkLogic {
	return &CreateLinkLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateLinkLogic) CreateLink(req *types.CreateLinkReq) (resp *types.CreateLinkResp, err error) {
	//从token中解析进ctx，ctx拿出userID
	userIdNum, ok := l.ctx.Value("user_id").(json.Number)
	if !ok {
		// l.Error()
		return nil, errors.New("JWT中缺少user_id")
	}

	userID, _ := userIdNum.Int64()
	// 调用rpc
	rpcResp, err := l.svcCtx.ShortenerRpc.Create(l.ctx, &shortener.CreateReq{
		Url:    req.Url,
		UserId: userID,
	})
	if err != nil {
		logx.Error()
		return nil, err
	}

	fullShortUrl := "localhost:8888/" + rpcResp.ShortUrl

	return &types.CreateLinkResp{
		ShortUrl: fullShortUrl,
	}, nil
}
