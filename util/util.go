package util

import (
	"crypto/md5"
	"fmt"
)

func MD5String(src []byte) string {
	return fmt.Sprintf("%x", md5.Sum(src))
}
