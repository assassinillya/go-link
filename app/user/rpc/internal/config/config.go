package config

import "github.com/zeromicro/go-zero/zrpc"

type Config struct {
	zrpc.RpcServerConf
	// 对应 yaml 里的 Mysql: -> DataSource:
	// json tag 里的名字必须和 yaml 里的 key 一模一样
	Mysql struct {
		DataSource string
	}
}
