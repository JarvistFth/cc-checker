package orgnization

import (
	"cc-checker/ccs/normal/smart-audit-publish/core/contract"
)

const (
	// 合约维护人员计数存储Key值前缀
	MaintainerPrefix = "maintainer-"

	// 合约维护人员计数存储Key值前缀
	MaintainerCountKey = "maintainer-count"
)

// 合约运维人员
type Maintainer struct {
	*Member
}

// 合约维护人员存储Key值
func (m *Maintainer) Key() string {
	return MaintainerPrefix + m.Member.Key()
}

// 合约维护人员存储计数Key值
func (m *Maintainer) CountKey() string {
	return MaintainerCountKey
}

// 合约维护人员存储当前计数
func (m *Maintainer) GetCount() uint32 {
	return m.ID + 1
}

// 从输入参数获取合约维护人员
func MaintainerFromString(args []string, index uint32,
	context contract.Context) (mt *Maintainer, err error) {
	var mem *Member
	if mem, err = MemberFromString(args, context); err != nil {
		return
	}

	mem.ID = index
	mt = &Maintainer{
		Member: mem,
	}
	return
}
