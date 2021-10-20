package invokes

import (
	"core/contract"
	"oracles/face/service/dummy"
	_ "oracles/identify/service"
	"strconv"
)

var validation contract.Validation = initValidation()

// 人脸识别规则验证
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

// 生成人脸识别规则实例
func initValidation() contract.Validation {
	// fixme 在真实商用环境下替换为完成好的service.FaceValidation
	//return &service.FaceValidation{}
	return &dummy.FaceValidation{}
}
