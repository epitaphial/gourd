package gourd

import (
	"crypto/md5"
	"encoding/hex"
	"time"
)

// 随机字符串生成
func getRandomString() string {
	hash := md5.New()
	hash.Write([]byte(string(time.Now().UnixNano()) + "curled"))
	return hex.EncodeToString(hash.Sum(nil))
}
