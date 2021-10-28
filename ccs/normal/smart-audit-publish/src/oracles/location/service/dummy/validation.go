package dummy

import (
	"cc-checker/core/dynamic/ccs/normal/smart-audit-publish/src/core/contract"
	"cc-checker/core/dynamic/ccs/normal/smart-audit-publish/src/oracles/location/service"
	"errors"
	"math"
	"strings"
)

const (
	// 地球半径
	earthRadius float64 = 6371000
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

type LocationValidation struct {
}

// 模拟验证位置规则
func (t *LocationValidation) Validate(id contract.ServiceRuleID, args []string) error {
	if len(args) < 1 {
		return errors.New("地理位置解析需要的参数不足")
	}
	positionArgs := strings.Split(args[0], " ")
	position, err := service.PositionFromString(positionArgs)
	if err != nil {
		return errors.New("地理位置解析出错: " + err.Error())
	}

	return t.dummyServiceValidate(position)
}

// 这里模拟在定位预言机服务上对是否处于有效地理范围内的验证操作
func (t *LocationValidation) dummyServiceValidate(position *service.Position) error {
	area, err := service.FromStrings(fixServiceValidationRule.Rules[0].Params)
	if err != nil {
		return err
	}

	// 将需要验证的地理位置和初始化地理验证范围中心点从经纬度转换为弧度之后，计算出以米为单位的两点距离
	dLat := t.toRadians(position.Lat - area.Center.Lat)
	dLon := t.toRadians(position.Lon - area.Center.Lon)
	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(t.toRadians(area.Center.Lat))*math.Cos(t.toRadians(position.Lat))*
			math.Sin(dLon/2)*math.Sin(dLon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	dist := earthRadius * c

	// 将计算出的两点之间距离与验证范围半径作比较，若超出半径长度则验证不成功
	if dist > area.Radius {
		return errors.New("地理位置超出正常工作范围")
	}
	return nil
}

// 将度转换为弧度工具方法
func (t *LocationValidation) toRadians(degrees float64) float64 {
	return degrees * (math.Pi / 180)
}

// 这里粗略初始化了一个北京的1千米以内的验证规则
func initValidationRule() *ServiceRules {
	return &ServiceRules{
		Rules: []RuleItem{
			{
				Logic:     contract.NONE,
				Condition: contract.IN,
				Params:    []string{"39.9", "116.3", "1000"},
			},
		},
	}
}
