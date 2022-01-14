package dummy

import (
	"cc-checker/ccs/normal/smart-audit-publish/core/contract"
)

type EntityIdentifyRegistration struct {
	dummyID uint32
}

// 模拟注册物体识别规则
func (f *EntityIdentifyRegistration) Register(args []string) (contract.ServiceRuleID, error) {
	// 由于假设是固定的验证规则，所以这里不需要额外工作
	rtn := contract.ServiceRuleID(f.dummyID)
	f.dummyID++
	return rtn, nil
}
