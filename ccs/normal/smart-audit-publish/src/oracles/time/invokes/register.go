package invokes

import (
	"cc-checker/core/dynamic/ccs/normal/smart-audit-publish/src/core/contract"
	"cc-checker/core/dynamic/ccs/normal/smart-audit-publish/src/oracles/time/service/dummy"
	"strconv"
)

var registration contract.Registration = initRegistration()

// 注册时间规则
func RegisterMain(args []string, context contract.Context) *contract.Response {
	id, err := registration.Register(args)
	if err != nil {
		return contract.Error("规则注册错误，详细信息：" + err.Error())
	}

	return &contract.Response{
		Payload: []byte(strconv.FormatUint(uint64(id), 32)),
	}
}

// 调用预言机注册时间规则
func initRegistration() contract.Registration {
	// fixme 在真实商用环境下替换为完成好的service.TimeRegistration
	//return &service.TimeRegistration{}
	return &dummy.TimeRegistration{}
}
