package utils

import (
	"crypto/md5"
	"encoding/hex"
)

func Hash(info string) string {
	abs := md5.Sum([]byte(info))
	return hex.EncodeToString(abs[:])
}
