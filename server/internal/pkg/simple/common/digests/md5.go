package digests

import (
	"crypto/md5"
	"encoding/hex"
)

func MD5(str string) string {
	return MD5Bytes([]byte(str))
}

func MD5Bytes(data []byte) string {
	h := md5.New()
	h.Write(data)
	return hex.EncodeToString(h.Sum(nil))
}
