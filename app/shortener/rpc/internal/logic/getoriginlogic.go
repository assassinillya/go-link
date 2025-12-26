package logic

import (
	"context"
	"errors"
	"fmt"
	"go-link/app/shortener/model"
	"gorm.io/gorm"

	"go-link/app/shortener/rpc/internal/svc"
	"go-link/app/shortener/rpc/shortener"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetOriginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetOriginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetOriginLogic {
	return &GetOriginLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetOriginLogic) GetOrigin(in *shortener.GetOriginReq) (*shortener.GetOriginResp, error) {
	// 先查redis
	cacheKey := fmt.Sprintf("short:%s", in.ShortCode)

	val, err := l.svcCtx.Redis.Get(context.Background(), cacheKey).Result()
	if err == nil {
		// 缓存命中，直接返回
		return &shortener.GetOriginResp{OriginUrl: val}, nil
	}

	// 查mysql
	var link model.Link
	err = l.svcCtx.DB.Where("short_code = ?", in.ShortCode).Take(&link).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("短链不存在")
		}
		l.Error("查询失败", err)
		return nil, errors.New("查询失败")
	}

	return &shortener.GetOriginResp{
		OriginUrl: link.OriginalURL,
	}, nil
}
