// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package config

import (
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	rest.RestConf // 包含 Host,Port等基础配置

	// 对应yaml里的 Auth
	Auth struct {
		AccessSecret string
		AccessExpire int64
	}

	// 对应 yaml 里的 UserRpc
	// 类型必须是 zrpc.RpcClientConf，这样才能解析 Etcd 的配置
	UserRpc      zrpc.RpcClientConf
	ShortenerRpc zrpc.RpcClientConf
}
