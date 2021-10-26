package dynamic

import (
	"container/list"
	"errors"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"strings"
)

type StubAgent struct {
	MockStub
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
	return s
}

func (stub *StubAgent) GetState(key string) ([]byte, error) {
	value := stub.State[key]
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

	log.Debug("MockStub", stub.Name, "Putting", key, value)
	stub.State[key] = value

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







// cc.invoke(stub[interface])
// mockstub <- stubInterface