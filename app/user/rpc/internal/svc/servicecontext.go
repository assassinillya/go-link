package svc

import (
	"go-link/app/user/rpc/internal/config"
	"go-link/pkg/db"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config config.Config
	DB     *gorm.DB
}

func NewServiceContext(c config.Config) *ServiceContext {

	//c.Mysql.DataSource就是yaml里填的
	db.InitMySQL(c.Mysql.DataSource)

	return &ServiceContext{
		Config: c,
		DB:     db.MySQL,
	}
}
