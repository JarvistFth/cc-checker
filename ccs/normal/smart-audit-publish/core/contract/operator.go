package contract

type LogicOperator string
type ConditionalOperator string

// 逻辑关系表达相关操作符
const (
	// 空规则
	NONE LogicOperator = "NONE"

	// 非, 目前仅支持单个规则的情况
	NOT LogicOperator = "NOT"

	// 并且
	AND LogicOperator = "AND"

	// 或
	OR LogicOperator = "OR"
)

// 条件表达式相关操作符
const (
	// 大于
	GT ConditionalOperator = ">"

	// 大于等于
	GE ConditionalOperator = ">="

	// 小于
	LT ConditionalOperator = "<"

	// 小于等于
	LE ConditionalOperator = "<="

	// 等于
	EQ ConditionalOperator = "=="

	// 在…范围内
	IN ConditionalOperator = "in"
)
