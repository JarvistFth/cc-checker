package dummy

import (
	"cc-checker/ccs/normal/smart-audit-publish/src/core/contract"
	"errors"
	"strconv"
	"time"
)

const (
	layout = "2006-01-02T15:04:05.000Z"
)

var (
	fixServiceValidationRule = initValidationRule()
)

// 这里模拟在预言机服务上的单个规则结构
type RuleItem struct {
	Logic     contract.LogicOperator
	Condition contract.ConditionalOperator
	Params    []string
}

// 这里模拟通过规则表达式解析出的一组规则结构
type ServiceRules struct {
	Rules []RuleItem
}

type TimeValidation struct {
}

// 模拟验证时间是否满足规则要求
func (t *TimeValidation) Validate(id contract.ServiceRuleID, args []string) error {
	if len(args) == 0 {
		return errors.New("不允许没有参数的验证")
	}

	value, err := time.Parse(layout, args[0])
	if err != nil {
		return errors.New("时间解析出错: " + err.Error())
	}

	return t.dummyServiceValidate(value)
}

// 验证时间，需要时间在时间规则的开始时间及结束时间之间，否则返回错误
func (t TimeValidation) dummyServiceValidate(value time.Time) error {
	startTime := t.getTimeFromRuleParams(fixServiceValidationRule.Rules[0].Params)
	endTime :=  t.getTimeFromRuleParams(fixServiceValidationRule.Rules[1].Params)

	if value.Before(startTime) || value.After(endTime) {
		return errors.New("时间超出正常工作范围")
	}
	return nil
}

// 从输入参数获取时间
func (t *TimeValidation) getTimeFromRuleParams(params []string) time.Time {
	hour, _ := strconv.ParseInt(params[0], 10, 32)
	now := time.Now()
	return time.Date(now.Year(), now.Month(), now.Day(),
		int(hour), 0, 0, 0, time.UTC)
}

// 这里初始化了一个"朝九晚六"的验证规则
func initValidationRule() *ServiceRules {
	return &ServiceRules{
		Rules: []RuleItem{
			{
				Logic:     contract.AND,
				Condition: contract.GE,
				Params:    []string{"9"},
			},
			{
				Logic:     contract.AND,
				Condition: contract.LE,
				Params:    []string{"18"},
			},
		},
	}
}
