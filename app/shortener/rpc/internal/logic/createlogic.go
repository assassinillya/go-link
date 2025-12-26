package logic

import (
	"context"
	"errors"
	"fmt"
	"go-link/app/shortener/model"
	"go-link/app/shortener/rpc/internal/svc"
	"go-link/app/shortener/rpc/shortener"
	"go-link/pkg/utils"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

// GlobalIDKey 这是一个全局ID计数器的Key,存在Redis里
const GlobalIDKey = "go_link:global_id"

// StartIDOffset 为了让生成的断连不至于太短（比如id=1变成“1”），我们给他一个初始的偏移量
// 10000000000在62进制中是"aUKYOA"，长度合适
const StartIDOffset = 10000000000

func NewCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateLogic {
	return &CreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateLogic) Create(in *shortener.CreateReq) (*shortener.CreateResp, error) {

	userID := in.UserId

	id, err := l.svcCtx.Redis.Incr(context.Background(), GlobalIDKey).Result()
	if err != nil {
		logx.Error("BizRedis INCR失败：", err)
		return nil, errors.New("内部服务错误")
	}

	globalID := id + StartIDOffset

	shortCode := utils.Base62Encode(globalID)

	link := model.Link{
		UserID:      uint(userID),
		OriginalURL: in.Url,
		ShortCode:   shortCode,
		VisitCount:  0,
	}

	if err = l.svcCtx.DB.Create(&link).Error; err != nil {
		logx.Error("mysql保存失败:", err)
		return nil, errors.New("保存失败")
	}

	cacheKey := fmt.Sprintf("short:%s", shortCode)

	err = l.svcCtx.Redis.Set(context.Background(), cacheKey, in.Url, 7*24*time.Hour).Err()
	if err != nil {
		logx.Error("Redis缓存写入失败: %v\n", err)
	}

	return &shortener.CreateResp{
		// 只返回了短码，没拼接
		ShortUrl: shortCode,
	}, nil
}
