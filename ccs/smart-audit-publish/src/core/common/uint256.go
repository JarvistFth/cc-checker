package common

// 定义Uint256字节数
const UINT256SIZE = 32

// 定义Uint256为32个字节
type Uint256 [UINT256SIZE]uint8

// 将Uint56转换为16进制字符串
func (u Uint256) String() string {
	return BytesToHexString(u.Bytes())
}

// 将Uint256转换为字节数组
func (u Uint256) Bytes() []byte {
	var x = make([]byte, UINT256SIZE)
	copy(x, u[:])
	return x
}
