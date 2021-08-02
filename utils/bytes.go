package utils

import (
	"crypto/md5"
	"encoding/hex"
	"unsafe"
)

// StringToBytes string to []byte
func StringToBytes(str string) []byte {
	return *(*[]byte)(unsafe.Pointer(&str))
}

// BytesToString []byte to string
func BytesToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func MD5String(val string) (result string) {
	b := md5.Sum(StringToBytes(val))
	return hex.EncodeToString(b[:])
}
