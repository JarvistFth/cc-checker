package invokes

import (
	"cc-checker/core/dynamic/ccs/normal/smart-audit-publish/src/core/contract"
	"cc-checker/core/dynamic/ccs/normal/smart-audit-publish/src/oracles/identify/service/dummy"
	"strconv"
)

var registration contract.Registration = initRegistration()

// 注册物体识别规则
func RegisterMain(args []string, context contract.Context) *contract.Response {
	id, err := registration.Register(args)
	if err != nil {
		return contract.Error("规则注册错误，详细信息：" + err.Error())
	}

	return &contract.Response{
		Payload: []byte(strconv.FormatUint(uint64(id), 32)),
	}
}

// 生成物体识别实例
func initRegistration() contract.Registration {
	// fixme 在真实商用环境下替换为完成好的service.EntityIdentifyRegistration
	//return &service.EntityIdentifyRegistration{}
	return &dummy.EntityIdentifyRegistration{}
}
