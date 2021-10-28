package service

import (
	"cc-checker/core/dynamic/ccs/normal/smart-audit-publish/src/core/contract"
	"errors"
)

type TimeValidation struct {
}

// 验证时间是否满足规则
func (t *TimeValidation) Validate(id contract.ServiceRuleID, args []string) error {
	if len(args) < 1 {
		return errors.New("验证规则所需参数不足")
	}

	return t.serviceValidate(id, args[0])
}

// 调用预言机验证时间是否满足规则
func (t *TimeValidation) serviceValidate(id contract.ServiceRuleID,
	valueExpression string) error {
	// fixme 实际商用时实现时间预言机，然后在这里调用预言机服务
	return nil
}
