package orgnization

import (
	"cc-checker/core/dynamic/ccs/normal/smart-audit-publish/src/core/contract"
	"cc-checker/core/dynamic/ccs/normal/smart-audit-publish/src/core/record"
)

const (
	// 审计当事人存储Key值前缀
	AuditeePrefix   = "auditee-"
	// 审计当事人计数存储Key值前缀
	AuditeeCountKey = "auditee-count"
)

// 审计人员
type Auditee struct {
	*Member
}

// 审计人员存储Key值
func (a *Auditee) Key() string {
	return AuditeePrefix + a.Member.Key()
}

// 审计人员计数Key值
func (a *Auditee) CountKey() string {
	return AuditeeCountKey
}

// 获取审计人员当前计数
func (a *Auditee) GetCount() uint32 {
	return a.ID + 1
}

// 从输入参数获取审计人员，ID根据当前存储的审计人员计数加1
func AuditeeFromString(args []string,
	context contract.Context) (auditee *Auditee, err error) {
	var mem *Member
	if mem, err = MemberFromString(args, context); err != nil {
		return
	}

	auditee = &Auditee{
		Member: mem,
	}
	auditee.ID, err = record.GetRecordCount(AuditeeCountKey, context)
	return
}
