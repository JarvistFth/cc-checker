package rules

import (
	"cc-checker/ccs/smart-audit-publish/src/core/contract"
	"fmt"
	"strconv"
)

type ValidationExpression struct {
	// 规则类型
	Type RuleType

	// 具体验证规则
	Expression string
}

func RegisterRules(expression []string, context contract.Context) (uint32, error) {
	op, expressions, err := Parse(expression)
	if err != nil {
		return 0, err
	}

	relation := &ValidationRelationship{
		Operator: op,
		Rules:    make(map[RuleType]contract.ServiceRuleID, 0),
	}
	for _, v := range expressions {
		ruleID, err := v.registerRule(context)
		if err != nil {
			return 0, err
		}

		relation.Rules[v.Type] = ruleID
	}

	return registerValidationRelationship(relation, context)
}

func (r *ValidationExpression) registerRule(
	context contract.Context) (contract.ServiceRuleID, error) {
	switch r.Type {
	case Time, Location, FaceRecognize, ObjectRecognize:
		return r.registerFromContract(string(r.Type), context)
	default:
		return 0, fmt.Errorf("编码为%d的类型尚未支持", r.Type)
	}
}

func (r *ValidationExpression) registerFromContract(contractName string,
	context contract.Context) (contract.ServiceRuleID, error) {
	args := []string{
		r.Expression,
	}

	rtn := context.InvokeContract(contractName, contract.RegisterFunctionName, args)
	if rtn.Err != nil {
		return 0, rtn.Err
	}

	id, err := strconv.Atoi(string(rtn.Payload))
	if err != nil {
		return 0, err
	}

	return contract.ServiceRuleID(id), nil
}
