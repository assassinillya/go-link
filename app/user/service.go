package user

import (
	"github.com/gin-gonic/gin"
	"go-link/app/user/model"
	"go-link/pkg/db"
	"go-link/pkg/utils"
	"net/http"
)

// RegisterReq 注册请求参数
type RegisterReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Register 用户注册接口
func Register(c *gin.Context) {
	var req RegisterReq
	// 这个函数可以将http请求的json数据解析并绑定到结构体
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}
	// 现在已经把用户请求传过来的数据变成了req结构体

	// 进行密码加密
	hashedPwd, err := utils.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "密码加密失败"})
		return
	}

	// 存入数据库
	user := model.User{
		Username: req.Username,
		Password: hashedPwd,
	}

	if result := db.MySQL.Create(&user); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "注册失败，用户名可能已存在"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"msg": "注册成功"})
}

// LoginReq 登录请求参数
type LoginReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func Login(c *gin.Context) {
	var req LoginReq
	// 将请求中的json与req进行绑定
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	var user model.User
	// 查询数据库中有没有这个用户
	// 在这一步中，user对象已经和数据库中对应的user绑定上了
	if err := db.MySQL.Where("username = ?", req.Username).Take(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "用户不存在"})
	}

	// 校验密码
	if !utils.CheckPassword(req.Password, user.Password) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "密码错误"})
	}

	// 生成JWT
	token, err := utils.GenerateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "token生成失败"})
	}

	// 返回token
	c.JSON(http.StatusOK, gin.H{
		"msg":   "登录成功",
		"token": token,
	})
}
