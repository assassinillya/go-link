package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-link/internal/gateway/middleware"
	"go-link/internal/shortener"
	"go-link/internal/user"
	"go-link/pkg/config"
	"go-link/pkg/db"
	"log"
)

func main() {
	// 加载配置
	config.InitConfig()

	// 初始化数据库
	db.InitMySQL(config.AppConfig.Database.MysqlDSN)

	// 初始化Redis
	db.InitRedis(
		config.AppConfig.Database.RedisAddr,
		config.AppConfig.Database.RedisPwd,
	)

	// 初始化 Gin 引擎
	r := gin.Default()

	// 注册一个简单的路由 （GET请求）
	// 当访问/ping时，返回JSON数据
	//r.GET("/ping", func(c *gin.Context) {
	//	c.JSON(http.StatusOK, gin.H{
	//		"message": "pong",
	//		"project": "go-link",
	//		"status":  "running",
	//	})
	//	// 测试一下能不能拿到 DB 对象
	//	if db.MySQL != nil {
	//		c.JSON(http.StatusOK, gin.H{"status": "Database Connected!"})
	//	} else {
	//		c.JSON(http.StatusOK, gin.H{"status": "Database Error"})
	//	}
	//
	//	if db.Redis != nil {
	//		c.JSON(http.StatusOK, gin.H{"status": "Redis Connected!"})
	//	} else {
	//		c.JSON(http.StatusOK, gin.H{"status": "Database Error"})
	//	}
	//})

	//userGroup := r.Group("/user")
	//{
	//	userGroup.POST("/register", user.Register)
	//	userGroup.POST("/login", user.Login)
	//}

	// 公开路由
	r.GET("/:code", shortener.Redirect)

	public := r.Group("/api/v1")
	{
		public.POST("/register", user.Register)
		public.POST("/login", user.Login)
	}

	// 鉴权路由
	// Use(middleware.JWTAuth()) 表示这个组下面的所有接口都要过中间件
	protected := r.Group("/api/v1")
	protected.Use(middleware.JWTAuth())
	{
		// 定义一个短链相关的路由组
		linkGroup := protected.Group("/link")
		{
			linkGroup.POST("/create", shortener.Create)
			linkGroup.GET("/list", shortener.List)
			linkGroup.PUT("/update", shortener.Update)
			linkGroup.DELETE("/delete", shortener.Delete)
		}
	}

	addr := fmt.Sprintf(":%s", config.AppConfig.Server.Port)
	err := r.Run(addr)
	if err != nil {
		log.Println("启动失败：", err)
		return
	}
}
