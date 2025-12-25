package utils

import "strings"

// 字符集：0-9, a-z, A-Z (共62个字符)
// 顺序打乱一点可以防止被轻易猜出规律，这里我们用标准顺序
const charset = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

// Base62Encode 十进制转62进制
func Base62Encode(num int64) string {
	if num == 0 {
		return string(charset[0])
	}

	var sb strings.Builder
	base := int64(len(charset))

	for num > 0 {
		rem := num % base
		sb.WriteByte(charset[rem])
		num /= base
	}

	return reverseString(sb.String())
}

// 翻转字符串
func reverseString(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

// Base62Decode 62进制转十进制
func Base62Decode(str string) int64 {
	var num int64
	base := int64(len(charset))

	for _, char := range str {
		index := strings.IndexRune(charset, char)
		num = num*base + int64(index)
	}
	return num
}
