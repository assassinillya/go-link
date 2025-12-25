package db

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

// GlobalDB 是一个全局的数据库对象，方便后续在其他包直接调用
// 注意：在实际大项目中通常使用依赖注入，但作为练手项目，全局变量是最简单的入手方式

var MySQL *gorm.DB

// InitMySQL 初始化 MySQL 连接
// dsn 格式: "user:password@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"

func InitMySQL(dsn string) {
	var err error
	MySQL, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("连接MySQL失败，原因: %v", err)
	} else {
		log.Println("MySQL连接成功")
	}
}
