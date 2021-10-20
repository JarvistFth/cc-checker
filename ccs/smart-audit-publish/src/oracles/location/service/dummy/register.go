package dummy

import "core/contract"

type LocationRegistration struct {
	dummyID uint32
}

// 模拟注册位置规则
func (t *LocationRegistration) Register(args []string) (contract.ServiceRuleID, error) {
	// 由于假设是固定的验证规则，所以这里不需要额外工作
	rtn := contract.ServiceRuleID(t.dummyID)
	t.dummyID++
	return rtn, nil
}
