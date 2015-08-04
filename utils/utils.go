package utils

import (
	"crypto/md5"
	"encoding/hex"
	"os"
)

// 返回文件的md5和大小
func GetMd5FromFile(filename string) (string, int) {
	file, err := os.Open(filename)
	if err != nil {
		return "", 0
	}
	defer file.Close()

	b := make([]byte, 8*1024)
	h := md5.New()
	size := 0

	for {
		n, err := file.Read(b)
		if err != nil {
			break
		}
		size += n
		h.Write(b[:n])
	}

	return hex.EncodeToString(h.Sum(nil)), size
}
