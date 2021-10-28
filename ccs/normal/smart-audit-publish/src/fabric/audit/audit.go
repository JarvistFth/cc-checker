package main

import (
	invokes2 "cc-checker/core/dynamic/ccs/normal/smart-audit-publish/src/core/invokes"
	contract2 "cc-checker/core/dynamic/ccs/normal/smart-audit-publish/src/fabric/contract"
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
	context := contract2.NewContext(stub)
	res := invokes2.InitMain(context)
	return contract2.Response(res)
}

// 审计合约方法调用
func (s *SmartAudit) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	context := contract2.NewContext(stub)
	args := context.GetArgs()

	switch context.GetFunctionName() {
	// 录入审计规则
	case invokes2.RegisterRules:
		return contract2.Response(invokes2.RegisterRulesMain(args, context))
	// 录入审计当事人
	case invokes2.RegisterAuditee:
		return contract2.Response(invokes2.RegisterAuditeeMain(args, context))
	// 录入项目
	case invokes2.RegisterProject:
		return contract2.Response(invokes2.RegisterProjectMain(args, context))
	// 录入审计事件
	case invokes2.AddEvent:
		return contract2.Response(invokes2.AddEventMain(args, context))
	// 根据审计当事人ID，查询审计当事人信息
	case invokes2.GetAuditee:
		return contract2.Response(invokes2.GetAuditeeMain(args, context))
	// 根据规则ID，查询规则信息
	case invokes2.GetRule:
		return contract2.Response(invokes2.GetRulesMain(args, context))
	// 根据项目ID，查询项目信息
	case invokes2.GetProject:
		return contract2.Response(invokes2.GetProjectMain(args, context))
	// 获取所有合约维护人员
	case invokes2.GetMaintainers:
		return contract2.Response(invokes2.GetMaintainersMain(context))
	// 获取所有审计事件
	case invokes2.QueryEvents:
		return contract2.Response(invokes2.QueryEventsMain(args, context))
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
