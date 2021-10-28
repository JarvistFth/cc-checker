package dynamic

import (
	"cc-checker/logger"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

var log = logger.GetLogger()

func CheckNonDitermined(cc shim.Chaincode, fn string, args []string) {
	b := MockParams(fn,args)
	var agents []*StubAgent
	for i:=0; i<10; i++{
		agent := NewStubAgent("cc",cc)
		agents = append(agents,agent)
		
		agent.MockInvoke(genUUID(),b)
	}

	if !ConflictMap.IsDetermined(){
		log.Warning("Non-determined risk detect here!!")
	}
}

func CheckOtherRules(cc shim.Chaincode, fn string, args []string) {
	b := MockParams(fn,args)
	agent := NewMockStub("cc", cc)
	agent.MockInvoke(genUUID(),b)
}

