package service

import (
	"core/contract"
	"errors"
)

type EntityIdentifyRegistration struct {
}

// 注册物体识别规则
func (f *EntityIdentifyRegistration) Register(args []string) (contract.ServiceRuleID, error) {
	if len(args) < 1 {
		return 0, errors.New("注册规则所需参数不足")
	}

	return f.serviceRegister(args[0])
}

// 调用语言机注册物体识别规则
func (f *EntityIdentifyRegistration) serviceRegister(
	expression string) (contract.ServiceRuleID, error) {
	// fixme 实际商用时实现物体识别预言机，然后在这里调用预言机服务
	return 0, nil
}
