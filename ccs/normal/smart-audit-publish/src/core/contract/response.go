package contract

import "errors"

// 用于定义一个合约的返回值
type Response struct {
	// 返回错误，如果没有错误则返回为空
	Err error

	// 具体的返回值
	Payload []byte
}

// 将Error转换为Response
func Error(msg string) *Response {
	return &Response{
		Err:     errors.New(msg),
		Payload: nil,
	}
}
