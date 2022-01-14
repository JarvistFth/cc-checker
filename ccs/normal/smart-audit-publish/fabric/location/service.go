package main

import (
	core "cc-checker/ccs/normal/smart-audit-publish/core/contract"
	"cc-checker/ccs/normal/smart-audit-publish/fabric/contract"
	"cc-checker/ccs/normal/smart-audit-publish/oracles/location/invokes"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
	"log"
)

// 用于处理与定位服务预言机交互的智能合约
type LocationService struct {
}

// 位置服务合约初始化
func (s *LocationService) Init(stub shim.ChaincodeStubInterface) peer.Response {
	// 初始化位置预言机服务相关信息……
	return shim.Success(nil)
}

// 位置服务合约方法调用
func (s *LocationService) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	context := contract.NewContext(stub)
	args := context.GetArgs()

	switch context.GetFunctionName() {
	// 注册位置规则
	case core.RegisterFunctionName:
		return contract.Response(invokes.RegisterMain(args, context))
	// 位置验证
	case core.ValidationFunctionName:
		return contract.Response(invokes.ValidateMain(args))
	default:
		return shim.Error(fmt.Sprintf("找不到名为%s的方法，调用失败",
			context.GetFunctionName()))
	}
}

// 位置服务合约主程序入口
func main() {
	if err := shim.Start(new(LocationService)); err != nil {
		log.Printf("智能合约启动出错，详细信息：%s", err)
	}
}
