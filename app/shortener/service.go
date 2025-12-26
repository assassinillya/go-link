package shortener

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go-link/app/shortener/model"
	"go-link/pkg/db"
	"go-link/pkg/utils"
	"log"
	"net/http"
	"time"
)

/*
要在高并发下快速生成唯一且短的字符串，业界通用的最强方案是 “发号器 + Base62 编码”。
发号器 (ID Generator)：我们需要一个绝对不重复的数字 ID。
普通做法：用 MySQL 自增 ID。缺点：高并发下数据库是瓶颈，还得先插入才能拿到 ID。
高并发做法：用 Redis 的 INCR 原子递增命令。Redis 单机能扛 10万+ QPS，完全满足需求。
Base62 编码：将 10 进制的数字 ID 转换成 62 进制（0-9, a-z, A-Z）。
比如 ID 100000000 -> 转换后可能是 6LAze。
这样就把很长的数字变成了一个很短的字符串。
*/

type CreateLinkReq struct {
	URL string `json:"url" binding:"required,url"` // binding:"url" 会自动校验是否是合法网址
}

// GlobalIDKey 这是一个全局ID计数器的Key,存在Redis里
const GlobalIDKey = "go_link:global_id"

// StartIDOffset 为了让生成的断连不至于太短（比如id=1变成“1”），我们给他一个初始的偏移量
// 10000000000在62进制中是"aUKYOA"，长度合适
const StartIDOffset = 10000000000

func Create(c *gin.Context) {
	// 获取用户ID
	userID, _ := c.Get("userID")

	// 解析参数
	var req CreateLinkReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请提供合法的URL"})
		return
	}

	ctx := context.Background()

	// 高并发核心
	// 使用Redis原子递增获取唯一ID
	id, err := db.Redis.Incr(ctx, GlobalIDKey).Result()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Redis获取原子ID出错"})
		return
	}

	// 加上偏移量，保证生成的短码长度好看一点
	globalID := id + StartIDOffset

	// Base62编码
	shortCode := utils.Base62Encode(globalID)

	// 转为数据库类型
	link := model.Link{
		UserID:      userID.(uint),
		OriginalURL: req.URL,
		ShortCode:   shortCode,
		VisitCount:  0,
	}

	// 存入
	if err := db.MySQL.Create(&link).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "保存失败"})
		log.Println(err)
		return
	}

	// 写入Redis缓存
	// Key:"short:xyz" Value:"http://Google.com"
	cacheKey := fmt.Sprintf("short:%s", shortCode)

	// 设置7表示7纳秒后过期（必须乘上时间单位），或者设置0表示永不过期
	err = db.Redis.Set(ctx, cacheKey, req.URL, 7*24*time.Hour).Err()
	if err != nil {
		fmt.Printf("Redis缓存写入失败: %v\n", err)
	}

	debug := db.Redis.Get(ctx, cacheKey)
	log.Println("存入的redis值:", debug)

	// 返回结果
	// 拼接完整的短链接地址给前端
	host := "http://localhost:8080/"
	c.JSON(http.StatusOK, gin.H{
		"short_url":  host + shortCode,
		"short_code": shortCode,
	})
}

func Redirect(c *gin.Context) {

	// 获取URL路径参数中的短码
	code := c.Param("code")

	cacheKey := fmt.Sprintf("short:%s", code)

	originalURL, err := db.Redis.Get(context.Background(), cacheKey).Result()

	// Redis 命中，直接跳转
	if err == nil {
		// 记录访问日志(异步往Redis里的一个List推送)
		go db.Redis.Incr(c, "visit_count:"+code)

		// 302 临时重定向（301是永久，会导致浏览器缓存，通常使用302）
		c.Redirect(http.StatusFound, originalURL)
	}

	// 接下来是Redis未命中，我们需要查询MySQL
	if !errors.Is(err, redis.Nil) {
		log.Println("Redis数据异常:", err)
	}

	var link model.Link

	result := db.MySQL.Where("short_code = ?", code).First(&link)
	// mysql也没查到
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "链接不存在或已过期"})
		log.Println(err)
	}

	// mysql查到之后需要回写缓存
	err = db.Redis.Set(c, cacheKey, link.OriginalURL, 7*24*time.Hour).Err()
	if err != nil {
		log.Println("会写redis失败:", err)
	}
	c.Redirect(http.StatusFound, originalURL)
}

/*
这是还未写逻辑的模拟请求处理
func Create(c *gin.Context) {
	// 从上下文中取出userID(中间件存进去的)
	// 因为 c.Get 返回的是 interface{}，需要断言成 uint
	value, exists := c.Get("userID") //.(uint) 这里不能这么做，首先不知道userID是否存在，其次，没有检查断言结果
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		return
	}

	userID, ok := value.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "未登录"})
		return
	}

	// 解析参数
	var req CreateLinkReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "URL参数不正确"})
		return
	}

	// 模拟返回
	c.JSON(http.StatusOK, gin.H{
		"msg":         "验证成功，准备生成短链",
		"user_id":     userID,
		"receive_url": req.URL,
	})
}
*/

type ListReq struct {
	Page     int `json:"page"`      // 页码
	PageSize int `json:"page_size"` // 每页数量
}

func List(c *gin.Context) {

	userID, _ := c.Get("userID")

	var req ListReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	log.Println(req.Page, req.PageSize)

	if req.Page <= 0 {
		req.Page = 1
	}

	if req.PageSize <= 0 || req.PageSize > 100 {
		req.PageSize = 10
	}

	var links []model.Link
	var total int64

	// 计算偏移量
	offset := (req.Page - 1) * req.PageSize

	db.MySQL.Model(&model.Link{}).Where("user_id = ?", userID).Count(&total)

	// 分页查询数据
	result := db.MySQL.Where("user_id = ?", userID).
		Order("created_at desc").
		Offset(offset).
		Limit(req.PageSize).
		Find(&links)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  links,
		"total": total,
		"page":  req.Page,
	})
}

type UpdateReq struct {
	ID          uint   `json:"id" binding:"required"`
	NewOriginal string `json:"new_url" binding:"required,url"`
}

func Update(c *gin.Context) {
	userID, _ := c.Get("userID")
	var req UpdateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	// 查询链接是否存在
	var link model.Link
	if err := db.MySQL.Where("id = ? AND user_id = ?", req.ID, userID).First(&link).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "链接不存在或无权修改"})
		return
	}

	// 更新数据库
	if err := db.MySQL.Model(&link).Update("original_url", req.NewOriginal).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新数据库失败"})
		return
	}

	// 删除redis中无效缓存
	cacheKey := fmt.Sprintf("short:%s", link.ShortCode)
	db.Redis.Del(context.Background(), cacheKey)

	c.JSON(http.StatusOK, gin.H{"msg": "更新成功"})
}

type DeleteReq struct {
	ID uint `json:"id" binding:"required"`
}

func Delete(c *gin.Context) {
	UserID, _ := c.Get("userID")
	var req DeleteReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	var link model.Link
	if err := db.MySQL.Where("id = ? AND user_id = ?", req.ID, UserID).First(&link).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "链接不存在或无权删除"})
	}

	// 数据库硬删除
	db.MySQL.Delete(&link)

	// 删除Redis缓存
	cacheKey := fmt.Sprintf("short:" + link.ShortCode)
	db.Redis.Del(context.Background(), cacheKey)

	c.JSON(http.StatusOK, gin.H{"msg": "删除成功"})
}
