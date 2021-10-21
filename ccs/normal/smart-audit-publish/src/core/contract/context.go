package contract

// 抽象出与具体区块链无关的抽象智能合约接口
type Context interface {
	// 获取合约被调用时的方法名
	GetFunctionName() string

	// 获取合约被调用时的参数
	GetArgs() []string

	// 存储一个键值对
	PutState(key string, value []byte) error

	// 根据键获取相应的值
	GetState(key string) ([]byte, error)

	// 通过键删除指定存储项
	DeleteState(key string) error

	// 通过指定起始和结束键来获取一个范围内的所有存储项
	GetStateByRange(startKey, endKey string) (Iterator, error)

	// 调用另一个智能合约，目前默认会在同一个链（或者channel）上部署所有合约
	InvokeContract(name, function string, args []string) Response
}

// 抽象一个用户获取存储项的迭代器
type Iterator interface {

	// 判断是否还有下一项，有则返回true，否则返回false
	HasNext() bool

	// 获取下一个存储项
	Next() (key string, value []byte, err error)

	// 关闭迭代器，并清理相关资源
	Close() error
}
