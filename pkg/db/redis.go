package db

import (
	"context"
	"github.com/redis/go-redis/v9"
	"log"
)

// Redis 依旧全局客户端
var Redis *redis.Client

func InitRedis(addr string, Password string) {
	Redis = redis.NewClient(&redis.Options{
		Addr:     addr,     // 例如 "127.0.0.1:6379"
		Password: Password, // 如果没有密码就传空字符串
		DB:       0,        // 默认使用 0 号数据库
	})

	ctx := context.Background()
	_, err := Redis.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("连接Redis失败，原因: %v", err)
	} else {
		log.Println("Redis 连接成功")
	}
}
