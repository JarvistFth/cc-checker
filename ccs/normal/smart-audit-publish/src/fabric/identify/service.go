package main

import (
	core "cc-checker/ccs/normal/smart-audit-publish/src/core/contract"
	contract2 "cc-checker/ccs/normal/smart-audit-publish/src/fabric/contract"
	invokes2 "cc-checker/ccs/normal/smart-audit-publish/src/oracles/identify/invokes"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
	"log"
)

// 用于处理与物体识别预言机交互的智能合约
type IdentifyService struct {
}

// 物体识别合约初始化
func (s *IdentifyService) Init(stub shim.ChaincodeStubInterface) peer.Response {
	// 初始化物体识别预言机服务相关信息……
	return shim.Success(nil)
}

// 物体识别合约方法调用
func (s *IdentifyService) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	context := contract2.NewContext(stub)
	args := context.GetArgs()

	switch context.GetFunctionName() {
	// 注册物体规则
	case core.RegisterFunctionName:
		return contract2.Response(invokes2.RegisterMain(args, context))
	// 物体验证
	case core.ValidationFunctionName:
		return contract2.Response(invokes2.ValidateMain(args))
	default:
		return shim.Error(fmt.Sprintf("找不到名为%s的方法，调用失败",
			context.GetFunctionName()))
	}
}

// 物体识别主程序入口
func main() {
	if err := shim.Start(new(IdentifyService)); err != nil {
		log.Printf("智能合约启动出错，详细信息：%s", err)
	}
}
