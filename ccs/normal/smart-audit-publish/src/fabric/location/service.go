package main

import (
	core "cc-checker/core/dynamic/ccs/normal/smart-audit-publish/src/core/contract"
	contract2 "cc-checker/core/dynamic/ccs/normal/smart-audit-publish/src/fabric/contract"
	invokes2 "cc-checker/core/dynamic/ccs/normal/smart-audit-publish/src/oracles/location/invokes"
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
	context := contract2.NewContext(stub)
	args := context.GetArgs()

	switch context.GetFunctionName() {
	// 注册位置规则
	case core.RegisterFunctionName:
		return contract2.Response(invokes2.RegisterMain(args, context))
	// 位置验证
	case core.ValidationFunctionName:
		return contract2.Response(invokes2.ValidateMain(args))
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
