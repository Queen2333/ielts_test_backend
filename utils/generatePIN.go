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
