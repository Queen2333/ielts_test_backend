package utils

import (
	"math/rand"
	"time"
)

// GeneratePIN 生成指定长度的随机字符串
func GeneratePIN(length int) string {
	rand.Seed(time.Now().UnixNano())

	pin := make([]byte, length)
	for i := 0; i < length; i++ {
		pin[i] = byte(rand.Intn(10) + '0')
	}

	return string(pin)
}

// 生成随机指定位数的数字
func GenerateRandomNumber(length int) string {
	const charset = "0123456789"
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}