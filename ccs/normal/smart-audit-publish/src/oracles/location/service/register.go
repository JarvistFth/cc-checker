package service

import (
	"cc-checker/core/dynamic/ccs/normal/smart-audit-publish/src/core/contract"
	"errors"
)

type LocationRegistration struct {
}

// 注册位置规则
func (t *LocationRegistration) Register(args []string) (contract.ServiceRuleID, error) {
	if len(args) < 1 {
		return 0, errors.New("注册规则所需参数不足")
	}

	return t.serviceRegister(args[0])
}

// 调用预言机注册位置规则
func (t *LocationRegistration) serviceRegister(
	expression string) (contract.ServiceRuleID, error) {
	// fixme 实际商用时实现定位预言机，然后在这里调用预言机服务
	return 0, nil
}
