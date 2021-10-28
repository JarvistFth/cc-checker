package contract

import (
	contract2 "cc-checker/core/dynamic/ccs/normal/smart-audit-publish/src/core/contract"
	"errors"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// 将Fabric平台上下文封装，并实现统一的数据库使用等方法接口，以支持不同链的调用
type ContextImpl struct {
	stub shim.ChaincodeStubInterface
}

// 获取输入参数中的方法名，参数中第一个字段为方法名
func (c *ContextImpl) GetFunctionName() string {
	fn, _ := c.stub.GetFunctionAndParameters()
	return fn
}

// 获取输入参数列表
func (c *ContextImpl) GetArgs() []string {
	_, args := c.stub.GetFunctionAndParameters()
	return args
}

// 存储数据
func (c *ContextImpl) PutState(key string, value []byte) error {
	return c.stub.PutState(key, value)
}

// 根据Key值查询之前存储的value值
func (c *ContextImpl) GetState(key string) ([]byte, error) {
	return c.stub.GetState(key)
}

// 删除Key值对应记录
func (c *ContextImpl) DeleteState(key string) error {
	return c.stub.DelState(key)
}

// 根据起始Key值及结束Key值，查询所有满足条件的值，返回迭代器
func (c *ContextImpl) GetStateByRange(startKey, endKey string) (contract2.Iterator, error) {
	raw, err := c.stub.GetStateByRange(startKey, endKey)
	if err != nil {
		return nil, err
	}
	return NewIterator(raw), nil
}

// 根据合约名、方法名及输入参数调用合约
func (c *ContextImpl) InvokeContract(name, function string, args []string) contract2.Response {
	fabricArgs := [][]byte{
		[]byte(function),
	}
	for _, v := range args {
		fabricArgs = append(fabricArgs, []byte(v))
	}
	// 这里channel参数为空则默认会发送到当前合约所在channel上
	rtn := c.stub.InvokeChaincode(name, fabricArgs, "")
	response := contract2.Response{
		Payload: rtn.Payload,
	}
	if rtn.Status != shim.OK {
		response.Err = errors.New(rtn.Message)
	}
	return response
}

// 生成ContextImpl实例
func NewContext(stub shim.ChaincodeStubInterface) *ContextImpl {
	return &ContextImpl{
		stub: stub,
	}
}
