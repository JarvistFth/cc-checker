package record

import (
	"cc-checker/ccs/normal/smart-audit-publish/src/core/contract"
)

// 用于抽象一条记录的信息
type Item interface {
	// Key值
	Key() string

	// 序列化存储除Key值外的所有数据
	Value() ([]byte, error)
}

// 存储一条记录，将接口中的Key、Value值一起存储到区块链中
func StoreItem(item Item, context contract.Context) error {
	value, err := item.Value()
	if err != nil {
		return err
	}
	return context.PutState(item.Key(), value)
}
