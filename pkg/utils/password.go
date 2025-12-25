package utils

import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) (string, error) {
	// cost是哈希算法的计算成本因子，它决定了哈希计算的复杂度和安全性级别
	// []byte表示转换成字节数组，还计算法在底层操作的是字节
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckPassword(password string, hash string) bool {
	// 传入用户请求中的密码 和 数据库中查到的加密后的密码 进行校验
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
