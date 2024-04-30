package number

import (
	"fmt"
	"github.com/satori/go.uuid"
	"math/rand"
	"strings"
	"time"
)

// GetRandNum 随机数
func GetRandNum(len int) int {
	if len <= 0 {
		return len
	}
	rand.Seed(time.Now().UnixNano())
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	num := r.Intn(len)

	return num
}

// GenerateTraceId 生成traceId
func GenerateTraceId() string {
	id := uuid.NewV4().String()
	return strings.Replace(id, "-", "", -1)
}

// GenerateSpanId 生成spanId
func GenerateSpanId(length int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	result := make([]byte, length)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := range result {
		result[i] = str[r.Intn(61)]
	}

	return fmt.Sprintf("%s", string(result))
}
