package rules

import (
	"cc-checker/ccs/smart-audit-publish/src/core/contract"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
)

// 验证规则结构, 用于定义在某一个预言机相关智能合约做具体验证的验证规则子项
type ValidationValue struct {
	// 规则类型
	Type RuleType

	// 验证值
	ActualValues string

	// 规则表达式在预言机中对应的ID
	ID contract.ServiceRuleID
}

// 验证一个规则表达式，验证时会解析出逻辑操作符以及每个验证规则子项，视情况具体验证
func ValidateRules(relationID uint32, expressions []string,
	context contract.Context) error {

	relation := &ValidationRelationship{
		Rules: make(map[RuleType]contract.ServiceRuleID, 0),
		ID:    relationID}
	value, err := context.GetState(relation.Key())
	if err != nil {
		return errors.New("规则ID对应的规则不存在，详细信息：" + err.Error())
	}

	if err := json.Unmarshal(value, relation); err != nil {
		return errors.New("规则解析失败，详细信息：" + err.Error())
	}

	items, err := parseRuleValues(expressions)
	if err != nil {
		return err
	}

	if err = setRuleIds(relation, items); err != nil {
		return err
	}

	// 根据逻辑操作符验证规则项
	switch relation.Operator {
	case contract.NONE:
		return nil
	case contract.AND:
		return validateAnd(items, context)
	case contract.OR:
		return validateOr(items, context)
	case contract.NOT:
		return validateNot(items, context)
	default:
		return errors.New("未找到有效地关系表达式操作符")
	}
}

// 当逻辑操作符为AND时，验证规则中包含的所有规则子项
func validateAnd(items []*ValidationValue, context contract.Context) error {
	for _, v := range items {
		if err := v.Validate(context); err != nil {
			return err
		}
	}
	return nil
}

// 当逻辑操作符为OR时，验证规则中包含的所有规则子项
func validateOr(items []*ValidationValue, context contract.Context) (
	lastErr error) {
	for _, v := range items {
		if err := v.Validate(context); err != nil {
			lastErr = err
		} else {
			return nil
		}
	}
	return
}

// 当逻辑操作符为NOT时，验证规则子项中的第一个值如果验证不通过则说明整个规则验证通过
func validateNot(items []*ValidationValue, context contract.Context) error {
	if len(items) == 0 {
		return errors.New("缺少规则验证项")
	}

	if err := items[0].Validate(context); err != nil {
		return nil
	} else {
		return errors.New("验证项不存在预期的错误")
	}
}

func setRuleIds(relation *ValidationRelationship, items []*ValidationValue) error {
	if len(items) != len(relation.Rules) {
		return errors.New("给定的规则值数与需要验证的规则数不匹配")
	}

	var ok bool
	for _, v := range items {
		v.ID, ok = relation.Rules[v.Type]
		if !ok {
			return fmt.Errorf("规则类型%s不在需要验证的列表中", string(v.Type))
		}
	}
	return nil
}

func parseRuleValues(expressions []string) ([]*ValidationValue, error) {
	var results []*ValidationValue
	for i := 0; i+1 < len(expressions); i += 2 {
		t, err := getRuleType(expressions[i])
		if err != nil {
			return nil, err
		}
		results = append(results, &ValidationValue{
			Type:         t,
			ActualValues: expressions[i+1],
		})
	}
	return results, nil
}

// 一个规则子项的验证逻辑，根据规则类型执行实际的调用逻辑
func (i *ValidationValue) Validate(context contract.Context) error {
	switch i.Type {
	case None:
		return nil
	case Time, Location, FaceRecognize, ObjectRecognize:
		return i.validateFromContract(string(i.Type), context)
	default:
		return fmt.Errorf("编码为%d的类型尚未支持", i.Type)
	}
}

// 实际调用预言机相关的智能合约以执行规则验证
func (i *ValidationValue) validateFromContract(contractName string,
	context contract.Context) error {
	args := []string{
		strconv.FormatUint(uint64(i.ID), 32),
		i.ActualValues,
	}

	if rtn := context.InvokeContract(contractName,
		contract.ValidationFunctionName, args); rtn.Err != nil {
		return rtn.Err
	}
	return nil
}
