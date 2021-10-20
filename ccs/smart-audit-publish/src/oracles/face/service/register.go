package service

import (
	"core/contract"
	"errors"
)

type FaceRegistration struct {
}

// 人脸识别规则注册
func (f *FaceRegistration) Register(args []string) (contract.ServiceRuleID, error) {
	if len(args) < 1 {
		return 0, errors.New("注册规则所需参数不足")
	}

	return f.serviceRegister(args[0])
}

// 调用预言机注册人脸识别规则
func (f *FaceRegistration) serviceRegister(
	expression string) (contract.ServiceRuleID, error) {
	// fixme 实际商用时实现人脸识别预言机，然后在这里调用预言机服务
	return 0, nil
}
