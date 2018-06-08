package library

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
)

//计算32字节的MD5值
func Mymd5(str string) string {
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(str))
	cipherStr := md5Ctx.Sum(nil)
	return hex.EncodeToString(cipherStr)
}

//计算sha256
func Mysha256(str string) string {
	sha256Ctx := sha256.New()
	sha256Ctx.Write([]byte(str))
	return hex.EncodeToString(sha256Ctx.Sum(nil))
}
