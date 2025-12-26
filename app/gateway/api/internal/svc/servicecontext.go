// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package svc

import (
	"github.com/zeromicro/go-zero/zrpc"
	"go-link/app/gateway/api/internal/config"
	"go-link/app/shortener/rpc/shortenerclient"
	"go-link/app/user/rpc/userclient"
)

type ServiceContext struct {
	Config       config.Config
	UserRpc      userclient.User
	ShortenerRpc shortenerclient.Shortener
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,

		// 初始化客户端
		// c.UserRpc 就是从 config.go -> yaml 里读出来的配置
		UserRpc:      userclient.NewUser(zrpc.MustNewClient(c.UserRpc)),
		ShortenerRpc: shortenerclient.NewShortener(zrpc.MustNewClient(c.ShortenerRpc)),
	}
}
