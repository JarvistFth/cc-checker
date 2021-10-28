package invokes

import (
	"cc-checker/core/dynamic/ccs/normal/smart-audit-publish/src/core/contract"
	"cc-checker/core/dynamic/ccs/normal/smart-audit-publish/src/oracles/time/service/dummy"
	"strconv"
)

var validation contract.Validation = initValidation()

// 验证时间是否满足规则
func ValidateMain(args []string) *contract.Response {
	if len(args) == 0 {
		return contract.Error("缺少规则ID")
	}

	value, err := strconv.ParseUint(args[0], 10, 32)
	if err != nil {
		return contract.Error("解析规则ID出错，详细信息：" + err.Error())
	}
	if err = validation.Validate(contract.ServiceRuleID(value), args[1:]); err != nil {
		return contract.Error("验证错误，详细信息：" + err.Error())
	}

	return &contract.Response{}
}

// 调用预言机验证时间是否满足规则
func initValidation() contract.Validation {
	// fixme 在真实商用环境下替换为完成好的service.TimeValidation
	//return &service.TimeValidation{}
	return &dummy.TimeValidation{}
}
