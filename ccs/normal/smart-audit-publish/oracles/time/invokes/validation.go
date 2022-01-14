package invokes

import (
	contract2 "cc-checker/ccs/normal/smart-audit-publish/core/contract"
	"cc-checker/ccs/normal/smart-audit-publish/oracles/time/service/dummy"
	"strconv"
)

var validation contract2.Validation = initValidation()

// 验证时间是否满足规则
func ValidateMain(args []string) *contract2.Response {
	if len(args) == 0 {
		return contract2.Error("缺少规则ID")
	}

	value, err := strconv.ParseUint(args[0], 10, 32)
	if err != nil {
		return contract2.Error("解析规则ID出错，详细信息：" + err.Error())
	}
	if err = validation.Validate(contract2.ServiceRuleID(value), args[1:]); err != nil {
		return contract2.Error("验证错误，详细信息：" + err.Error())
	}

	return &contract2.Response{}
}

// 调用预言机验证时间是否满足规则
func initValidation() contract2.Validation {
	// fixme 在真实商用环境下替换为完成好的service.TimeValidation
	//return &service.TimeValidation{}
	return &dummy.TimeValidation{}
}
