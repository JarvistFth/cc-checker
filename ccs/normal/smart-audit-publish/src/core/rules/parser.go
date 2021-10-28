package rules

import (
	"cc-checker/ccs/normal/smart-audit-publish/src/core/contract"
	"errors"
)

// fixme 这里只做简单的一层（直接函数调用）或者两层（带逻辑操作符的函数组合调用）的情况，复杂的嵌套调用待后续实现

func Parse(words []string) (contract.LogicOperator, []ValidationExpression, error) {
	op, err := getLogicOperator(words[0])
	if err != nil {
		return op, nil, err
	}

	expressions, err := parseRules(op, words[1:])
	return op, expressions, err
}

func parseRules(op contract.LogicOperator, words []string) (
	expressions []ValidationExpression, err error) {
	// 为空规则时参数视为空
	if op == contract.NONE {
		return
	}

	if len(words) < 2 {
		err = errors.New("规则解析参数不足")
		return
	}

	var expr ValidationExpression
	for i := 0; i+1 < len(words); i += 2 {
		expr, err = parseSingleRule(words[i], words[i+1])
		if err != nil {
			return
		}
		expressions = append(expressions, expr)

		// 为非操作符时，解析完第一个规则即返回
		if op == contract.NOT {
			return
		}
	}
	return
}

func parseSingleRule(t string, expr string) (ValidationExpression, error) {
	ruleType, err := getRuleType(t)
	if err != nil {
		return ValidationExpression{}, err
	}

	return ValidationExpression{
		Type:       ruleType,
		Expression: expr,
	}, nil
}

// 获取最前面的逻辑操作符
func getLogicOperator(word string) (contract.LogicOperator, error) {
	switch word {
	case string(contract.NONE), string(contract.NOT), string(contract.AND),
		string(contract.OR):
		return contract.LogicOperator(word), nil
	default:
		return "", errors.New("未找到逻辑操作符")
	}
}

// 获取最前面的规则类型
func getRuleType(word string) (RuleType, error) {
	switch word {
	case string(Time):
		return Time, nil
	case string(Location):
		return Location, nil
	case string(FaceRecognize):
		return FaceRecognize, nil
	case string(ObjectRecognize):
		return ObjectRecognize, nil
	default:
		return None, errors.New("未找到规则")
	}
}
