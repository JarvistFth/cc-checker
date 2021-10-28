package main

import (
	"cc-checker/ccs/normal/smart-audit-publish/src/core/invokes"
	"cc-checker/ccs/normal/smart-audit-publish/src/fabric/contract"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
	"log"
)

// Fabric平台上的审计业务智能合约实现
type SmartAudit struct {
}

// 审计合约初始化
func (s *SmartAudit) Init(stub shim.ChaincodeStubInterface) peer.Response {
	context := contract.NewContext(stub)
	res := invokes.InitMain(context)
	return contract.Response(res)
}

// 审计合约方法调用
func (s *SmartAudit) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	context := contract.NewContext(stub)
	args := context.GetArgs()

	switch context.GetFunctionName() {
	// 录入审计规则
	case invokes.RegisterRules:
		return contract.Response(invokes.RegisterRulesMain(args, context))
	// 录入审计当事人
	case invokes.RegisterAuditee:
		return contract.Response(invokes.RegisterAuditeeMain(args, context))
	// 录入项目
	case invokes.RegisterProject:
		return contract.Response(invokes.RegisterProjectMain(args, context))
	// 录入审计事件
	case invokes.AddEvent:
		return contract.Response(invokes.AddEventMain(args, context))
	// 根据审计当事人ID，查询审计当事人信息
	case invokes.GetAuditee:
		return contract.Response(invokes.GetAuditeeMain(args, context))
	// 根据规则ID，查询规则信息
	case invokes.GetRule:
		return contract.Response(invokes.GetRulesMain(args, context))
	// 根据项目ID，查询项目信息
	case invokes.GetProject:
		return contract.Response(invokes.GetProjectMain(args, context))
	// 获取所有合约维护人员
	case invokes.GetMaintainers:
		return contract.Response(invokes.GetMaintainersMain(context))
	// 获取所有审计事件
	case invokes.QueryEvents:
		return contract.Response(invokes.QueryEventsMain(args, context))
	// 其它不支持方法调用则返回错误
	default:
		return shim.Error(fmt.Sprintf("找不到名为%s的方法，调用失败",
			context.GetFunctionName()))
	}
}

// 审计合约主程序入口
func main() {
	if err := shim.Start(new(SmartAudit)); err != nil {
		log.Printf("智能合约启动出错，详细信息：%s", err)
	}
}
