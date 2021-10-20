package dummy

import "core/contract"

type FaceRegistration struct {
	dummyID uint32
}

// 模拟人脸识别规则注册
func (f *FaceRegistration) Register(args []string) (contract.ServiceRuleID, error) {
	// 由于假设是固定的验证规则，所以这里不需要额外工作
	rtn := contract.ServiceRuleID(f.dummyID)
	f.dummyID++
	return rtn, nil
}
