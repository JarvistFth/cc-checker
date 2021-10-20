package invokes

const (
	// 录入审计规则，合约调用方法名
	RegisterRules = "registerRules"
	// 录入审计当事人，合约调用方法名
	RegisterAuditee = "registerAuditee"
	// 录入项目，合约调用方法名
	RegisterProject = "registerProject"
	// 录入审计事件，合约调用方法名
	AddEvent = "addEvent"

	// 查询审计规则，合约调用方法名
	GetRule = "getRule"
	// 查询审计维护人员，合约调用方法名
	GetMaintainers = "getMaintainers"
	// 查询审计当事人，合约调用方法名
	GetAuditee = "getAuditee"
	// 查询项目，合约调用方法名
	GetProject = "getProject"
	// 查询审计事件，合约调用方法名
	QueryEvents = "queryEvents"
)
