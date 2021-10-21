package invokes

import (
	contract2 "cc-checker/ccs/normal/smart-audit-publish/src/core/contract"
	rules2 "cc-checker/ccs/normal/smart-audit-publish/src/core/rules"
	"fmt"
	"log"
	"strconv"
)

// 注册规则，返回规则ID
func RegisterRulesMain(args []string, context contract2.Context) *contract2.Response {
	ruleID, err := rules2.RegisterRules(args, context)
	if err != nil {
		log.Println("注册规则失败：", err.Error())
		return contract2.Error(fmt.Sprint("注册规则失败，详细信息：", err))
	}

	log.Println("审计规则录入成功，规则ID：", ruleID)
	return &contract2.Response{
		Payload: []byte(strconv.FormatUint(uint64(ruleID), 32)),
	}
}

// 根据规则ID获取规则信息
func GetRulesMain(args []string, context contract2.Context) *contract2.Response {
	if len(args) == 0 {
		return contract2.Error("查询失败，需要提供规则ID")
	}

	ruleID, err := strconv.ParseUint(args[0], 10, 32)
	if err != nil {
		return contract2.Error(fmt.Sprintf("解析规则ID出错，详细信息：%s", err.Error()))
	}

	rule := rules2.ValidationRelationship{
		Rules: make(map[rules2.RuleType]contract2.ServiceRuleID, 0),
		ID:    uint32(ruleID)}
	ruleBuf, err := context.GetState(rule.Key())
	if err != nil {
		return contract2.Error(fmt.Sprintf("获取规则出错，详细信息：%s", err.Error()))
	}

	return &contract2.Response{Payload: ruleBuf}
}
