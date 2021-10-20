package rules

import (
	"cc-checker/ccs/smart-audit-publish/src/core/contract"
	"cc-checker/ccs/smart-audit-publish/src/core/record"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
)

const (
	rulePrefix   = "rule-"
	ruleCountKey = "rule-count"
)

type ValidationRelationship struct {
	// 逻辑操作符
	Operator contract.LogicOperator

	// 用于记录一组规则，Key值对应一个规则类型，Value值为注册规则表达式时预言机服务返回的相应的ID值
	Rules map[RuleType]contract.ServiceRuleID

	// 规则唯一标识
	ID uint32
}

func (v *ValidationRelationship) CountKey() string {
	return ruleCountKey
}

func (v *ValidationRelationship) GetCount() uint32 {
	return v.ID + 1
}

func (v *ValidationRelationship) Key() string {
	return rulePrefix + strconv.FormatUint(uint64(v.ID), 10)
}

func (v *ValidationRelationship) Value() ([]byte, error) {
	value, err := json.Marshal(v)
	if err != nil {
		return nil, errors.New("failed to marshal validation relationship")
	}

	return value, nil
}

func registerValidationRelationship(p *ValidationRelationship,
	context contract.Context) (uint32, error) {
	count, err := record.GetRecordCount(ruleCountKey, context)
	if err != nil {
		return 0, err
	}
	p.ID = count

	if err := record.StoreItem(p, context); err != nil {
		return 0, fmt.Errorf("审计业务%s存储失败，详细信息：%s", p.Key(), err)
	}
	if err := record.StoreCount(p, context); err != nil {
		return 0, fmt.Errorf("审计业务%s相应的索引值存储失败，详细信息：%s",
			p.Key(), err)
	}
	return p.ID, nil
}
