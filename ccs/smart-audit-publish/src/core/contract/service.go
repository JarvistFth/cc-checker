package contract

const (
	ValidationFunctionName = "validation"
	RegisterFunctionName = "register"
)

// 用于定义服务中的规则ID
type ServiceRuleID uint32

// 定义Validation接口
type Validation interface {
	Validate(id ServiceRuleID, args []string) error
}

// 定义Registration接口
type Registration interface {
	Register(args []string) (ServiceRuleID, error)
}
