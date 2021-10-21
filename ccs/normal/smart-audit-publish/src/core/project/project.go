package project

import (
	"encoding/json"
	"errors"
	"strconv"
)

const (
	// 项目存储Key值前缀
	Prefix   = "project-"
	// 项目计数存储Key值前缀
	CountKey = "project-count"
)

// 用于定义一个审计业务的结构
type Project struct {
	// 业务名称
	Name string

	// 业务唯一标识
	ID uint32

	// 业务相关描述
	Description string

	// 项目审计当事人及规则
	AuditeeRulesMap map[string]string
}

// 项目存储计数的Key值
func (p *Project) CountKey() string {
	return CountKey
}

// 获取项目存储计数
func (p *Project) GetCount() uint32 {
	return p.ID + 1
}

// 项目存储Key值
func (p *Project) Key() string {
	return Prefix + strconv.FormatUint(uint64(p.ID), 10)
}

// 项目存储Value值
func (p *Project) Value() ([]byte, error) {
	value, err := json.Marshal(p)
	if err != nil {
		return nil, errors.New("failed to marshal project")
	}

	return value, nil
}
