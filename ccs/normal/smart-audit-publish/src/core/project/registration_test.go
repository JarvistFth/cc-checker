package project

import (
	"core/common"
	"core/contract"
	"core/orgnization"
	"core/rules"
	"log"
	"testing"
)

func TestRegistration_Value(t *testing.T) {
	testMap := make(map[rules.RuleType]contract.ServiceRuleID, 0)
	testMap[rules.Time] = 123
	testMap[rules.Location] = 345

	registration := Registration{
		AuditeeSpecification: AuditeeSpecification{
			ID: common.Uint256{1, 2, 3},
			Auditee: orgnization.Auditee{
				Member: &orgnization.Member{
					Name: "member a",
					ID:   1,
				},
			},
			Project: Project{
				Name:        "project a",
				ID:          2,
				Description: "it is for test",
			},
			Rule: rules.ValidationRelationship{
				Operator: "test",
				Rules:    testMap,
				ID:       3,
			},
		},
		Timestamp: 1232321230,
		Params:    []string{"time", ">=9 <=18"},
		Index:     4,
	}
	result, err := registration.Value()
	if err != nil {
		log.Println(err)
	}
	log.Println(string(result))
}
