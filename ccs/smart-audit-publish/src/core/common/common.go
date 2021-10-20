package common

import (
	"encoding/hex"
)

// 将字节数组转换为16进制字符串
func BytesToHexString(data []byte) string {
	return hex.EncodeToString(data)
}

// 将uint32转为字节数组
func Uint32ToBytes(n uint32) []byte {
	return []byte{
		byte(n),
		byte(n >> 8),
		byte(n >> 16),
		byte(n >> 24),
	}
}
