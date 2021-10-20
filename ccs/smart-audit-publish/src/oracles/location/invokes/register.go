package invokes

import (
	"core/contract"
	_ "oracles/location/service"
	"oracles/location/service/dummy"
	"strconv"
)

var registration contract.Registration = initRegistration()

// 注册位置规则
func RegisterMain(args []string, context contract.Context) *contract.Response {
	id, err := registration.Register(args)
	if err != nil {
		return contract.Error("规则注册错误，详细信息：" + err.Error())
	}

	return &contract.Response{
		Payload: []byte(strconv.FormatUint(uint64(id), 32)),
	}
}

// 生成位置服务实例
func initRegistration() contract.Registration {
	// fixme 在真实商用环境下替换为完成好的service.LocationRegistration
	//return &service.LocationRegistration{}
	return &dummy.LocationRegistration{}
}
