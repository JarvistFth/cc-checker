package dynamic

import (
	"container/list"
	"errors"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"reflect"
	"strings"
)


type StubStates map[*StubAgent]map[string][]byte

var ConflictMap StubStates
var StubAgents []*StubAgent

func (s StubStates) IsDetermined() bool{
	for i,stubi := range StubAgents{
		for j,stubj := range StubAgents{
			if i != j{
				if !reflect.DeepEqual(s[stubi],s[stubj]){
					return false
				}
			}
		}
	}
	return true
}

type StubAgent struct {
	 MockStub
	 RWSet map[string]bool
}

func init() {
	ConflictMap = make(map[*StubAgent]map[string][]byte)
}

func NewStubAgent(name string, cc shim.Chaincode) *StubAgent {
	log.Debug("MockStub(", name, cc, ")")
	s := new(StubAgent)
	s.Name = name
	s.CC = cc
	s.State = make(map[string][]byte)
	s.PvtState = make(map[string]map[string][]byte)
	s.EndorsementPolicies = make(map[string]map[string][]byte)
	s.Invokables = make(map[string]*MockStub)
	s.Keys = list.New()
	s.ChaincodeEventsChannel = make(chan *pb.ChaincodeEvent, 100) //define large capacity for non-blocking setEvent calls.
	s.Decorations = make(map[string][]byte)
	//if ConflictMap == nil{
	//	ConflictMap = make(map[string]string)
	//}
	s.RWSet = make(map[string]bool)
	StubAgents = append(StubAgents,s)
	return s
}

func (stub *StubAgent) GetState(key string) ([]byte, error) {
	value := stub.State[key]
	if _, ok := stub.RWSet[key]; ok{
		log.Warning("invoke has Read-Your-Write error!!")
	}
	log.Debug("MockStub", stub.Name, "Getting", key, value)
	fmt.Println("hello getstate!")
	return value, nil
}

func (stub *StubAgent) PutState(key string, value []byte) error {
	if stub.TxID == "" {
		err := errors.New("cannot PutState without a transactions - call stub.MockTransactionStart()?")
		log.Errorf("%+v", err)
		return err
	}

	// If the value is nil or empty, delete the key
	if len(value) == 0 {
		log.Debug("MockStub", stub.Name, "PutState called, but value is nil or empty. Delete ", key)
		return stub.DelState(key)
	}

	log.Debug("MockStub", stub.Name, "Putting", key, string(value))
	stub.State[key] = value
	stub.RWSet[key] = true
	ConflictMap[stub] = stub.State



	// insert key into ordered list of keys
	for elem := stub.Keys.Front(); elem != nil; elem = elem.Next() {
		elemValue := elem.Value.(string)
		comp := strings.Compare(key, elemValue)
		log.Debug("MockStub", stub.Name, "Compared", key, elemValue, " and got ", comp)
		if comp < 0 {
			// key < elem, insert it before elem
			stub.Keys.InsertBefore(key, elem)
			log.Debug("MockStub", stub.Name, "Key", key, " inserted before", elem.Value)
			break
		} else if comp == 0 {
			// keys exists, no need to change
			log.Debug("MockStub", stub.Name, "Key", key, "already in State")
			break
		} else { // comp > 0
			// key > elem, keep looking unless this is the end of the list
			if elem.Next() == nil {
				stub.Keys.PushBack(key)
				log.Debug("MockStub", stub.Name, "Key", key, "appended")
				break
			}
		}
	}

	// special case for empty Keys list
	if stub.Keys.Len() == 0 {
		stub.Keys.PushFront(key)
		log.Debug("MockStub", stub.Name, "Key", key, "is first element in list")
	}

	return nil
}

func (stub *StubAgent) MockInit(uuid string, args [][]byte) pb.Response {
	stub.Args = args
	stub.MockTransactionStart(uuid)
	res := stub.CC.Init(stub)
	stub.MockTransactionEnd(uuid)
	return res
}

func (stub *StubAgent) InvokeChaincode(chaincodeName string, args [][]byte, channel string) pb.Response {

	log.Warning("call cross_channel_invoke!!")
	if channel != "" {
		chaincodeName = chaincodeName + "/" + channel
	}
	// TODO "Args" here should possibly be a serialized pb.ChaincodeInput
	otherStub := stub.Invokables[chaincodeName]
	log.Debug("MockStub", stub.Name, "Invoking peer chaincode", otherStub.Name, args)
	//	function, strings := getFuncArgs(Args)
	res := otherStub.MockInvoke(stub.TxID, args)
	log.Debug("MockStub", stub.Name, "Invoked peer chaincode", otherStub.Name, "got", fmt.Sprintf("%+v", res))
	return res
}

// Invoke this chaincode, also starts and ends a transaction.
func (stub *StubAgent) MockInvoke(uuid string, args [][]byte) pb.Response {
	stub.Args = args
	stub.MockTransactionStart(uuid)
	res := stub.CC.Invoke(stub)
	stub.MockTransactionEnd(uuid)
	return res


}






// cc.invoke(stub[interface])
// mockstub <- stubInterface