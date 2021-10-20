package project

import (
	"encoding/json"
	"errors"
	"strconv"
)

const (
	// 审计事件计数存储Key值后缀
	RegistrationCountKeySuffix = "_count"
)

type RegistrationCount map[string]uint32

// 用于记录一次信息登录
type Registration struct {
	// 审计当事人规范对象
	AuditeeSpecification

	// 登录发生的时间
	Timestamp uint64

	// 登录涉及的参数
	Params []string

	// 用于标记隶属于一个业务下该审计当事人第几次录入
	Index uint32
}

// 审计事件计数存储Key值
func (r *Registration) CountKey() string {
	return r.ID.String() + RegistrationCountKeySuffix
}

// 获取审计事件计数
func (r *Registration) GetCount() uint32 {
	return r.Index + 1
}

// 审计事件存储的Key值
func (r *Registration) Key() string {
	return r.ID.String() + strconv.FormatUint(uint64(r.Index), 10)
}

// 审计事件存储的Value值
func (r *Registration) Value() ([]byte, error) {
	value, err := json.Marshal(r)
	if err != nil {
		return nil, errors.New("failed to marshal auditee specification")
	}

	return value, nil
}
