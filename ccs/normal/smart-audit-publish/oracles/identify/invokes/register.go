package invokes

import (
	contract2 "cc-checker/ccs/normal/smart-audit-publish/core/contract"
	"cc-checker/ccs/normal/smart-audit-publish/oracles/identify/service/dummy"
	"strconv"
)

var registration contract2.Registration = initRegistration()

// 注册物体识别规则
func RegisterMain(args []string, context contract2.Context) *contract2.Response {
	id, err := registration.Register(args)
	if err != nil {
		return contract2.Error("规则注册错误，详细信息：" + err.Error())
	}

	return &contract2.Response{
		Payload: []byte(strconv.FormatUint(uint64(id), 32)),
	}
}

// 生成物体识别实例
func initRegistration() contract2.Registration {
	// fixme 在真实商用环境下替换为完成好的service.EntityIdentifyRegistration
	//return &service.EntityIdentifyRegistration{}
	return &dummy.EntityIdentifyRegistration{}
}
