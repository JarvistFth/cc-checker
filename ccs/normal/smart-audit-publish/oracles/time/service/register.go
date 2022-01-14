package service

import (
	"cc-checker/ccs/normal/smart-audit-publish/core/contract"
	"errors"
)

type TimeRegistration struct {
}

// 注册时间规则
func (t *TimeRegistration) Register(args []string) (contract.ServiceRuleID, error) {
	if len(args) < 1 {
		return 0, errors.New("注册规则所需参数不足")
	}

	return t.serviceRegister(args[0])
}

// 调用预言机注册时间规则
func (t *TimeRegistration) serviceRegister(
	expression string) (contract.ServiceRuleID, error) {
	// fixme 实际商用时实现时间预言机，然后在这里调用预言机服务
	return 0, nil
}
