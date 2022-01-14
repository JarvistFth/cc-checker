package service

import (
	"cc-checker/ccs/normal/smart-audit-publish/core/contract"
	"errors"
)

type LocationValidation struct {
}

// 位置规则验证
func (t *LocationValidation) Validate(id contract.ServiceRuleID, args []string) error {
	if len(args) < 1 {
		return errors.New("验证规则所需参数不足")
	}

	return t.serviceValidate(id, args[0])
}

// 调用预言机验证位置是否满足规则
func (t *LocationValidation) serviceValidate(id contract.ServiceRuleID,
	valueExpression string) error {
	// fixme 实际商用时实现定位预言机，然后在这里调用预言机服务
	return nil
}
