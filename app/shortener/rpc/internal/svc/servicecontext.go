package svc

import (
	"github.com/redis/go-redis/v9"
	"go-link/app/shortener/rpc/internal/config"
	"go-link/pkg/db"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config config.Config
	DB     *gorm.DB
	Redis  *redis.Client
}

func NewServiceContext(c config.Config) *ServiceContext {
	// 初始化mysql redis
	db.InitRedis(c.BizRedis.Host, c.BizRedis.Pass)
	db.InitMySQL(c.Mysql.DataSource)
	return &ServiceContext{
		Config: c,
		DB:     db.MySQL, // 全局变量
		Redis:  db.Redis, // 全局变量
	}
}
