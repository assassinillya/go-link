package config

import (
	"github.com/spf13/viper"
	"log"
)

var AppConfig Config

type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	JWT      JWTConfig      `mapstructure:"jwt"`
}

type ServerConfig struct {
	Port string `mapstructure:"port"`
}

type DatabaseConfig struct {
	MysqlDSN  string `mapstructure:"mysql_dsn"`
	RedisAddr string `mapstructure:"redis_addr"`
	RedisPwd  string `mapstructure:"redis_pwd"`
}

type JWTConfig struct {
	Secret string `mapstructure:"secret"`
}

func InitConfig() {
	viper.SetConfigName("config") // 配置文件名 (不带后缀)
	viper.SetConfigType("yaml")   // 文件类型
	viper.AddConfigPath(".")      // 查找路径：当前目录

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("读取配置文件失败: %v", err)
	}

	// 把读取到的配置反序列化到结构体中
	if err := viper.Unmarshal(&AppConfig); err != nil {
		log.Fatalf("解析配置文件失败: %v", err)
	}

	log.Println("配置加载成功")
}
