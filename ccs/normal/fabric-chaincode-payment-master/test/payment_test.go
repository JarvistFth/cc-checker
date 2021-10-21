package main_test

import (
	"github.com/c12msr/fabrictest/main"
	"testing"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"fmt"
)




//Check if chaincode initialized correctly
func checkInit(t *testing.T, stub *shim.MockStub, args [][]byte) {
	res := stub.MockInit("1", args)
	if res.Status != shim.OK {
		fmt.Println("Init failed", string(res.Message))
		t.FailNow()
	}
}

func checkState(t *testing.T, stub *shim.MockStub, args [][]byte) {
	name := string(args[0])
	bytes := stub.State[name]

	if bytes == nil {
		fmt.Println("State", name, "failed to get value")
		t.FailNow()
	}
	value := args[1]
	if string(bytes) != string(value) {
		fmt.Println("State value", name, "was not", "as expected")
		t.FailNow()
	}
}

func checkQuery(t *testing.T, stub *shim.MockStub, args [][]byte, value string) {
	res := stub.MockInvoke("1", args)

	name := string(args[1])
	if res.Status != shim.OK {
		fmt.Println("Query", name, "failed", string(res.Message))
		t.FailNow()
	}
	if res.Payload == nil {
		fmt.Println("Query", name, "failed to get value")
		t.FailNow()
	}
	if string(res.Payload) != value {
		fmt.Println("Query value", name, "was not", value, "as expected")
		t.FailNow()
	}
}

func checkInvoke(t *testing.T, stub *shim.MockStub, args [][]byte) {
	res := stub.MockInvoke("1", args);
	if res.Status != shim.OK {
		fmt.Println("Invoke", args, "failed", string(res.Message))
		t.FailNow()
	}
}

//Test if init works as intended
func TestChaincode_Init(t *testing.T) {
	tcc := new(main.SimpleChaincode)
	stub := shim.NewMockStub("test", tcc)
	// Init Member1=100 Member2=200
	checkInit(t, stub, [][]byte{[]byte("init"), []byte("Alice"), []byte("100"), []byte("Bob"), []byte("200")})

	checkState(t, stub, [][]byte{[]byte("Alice"), []byte("100")})
	checkState(t, stub, [][]byte{[]byte("Bob"), []byte("200")})
}

func TestChaincode_Query(t *testing.T) {
	tcc := new(main.SimpleChaincode)
	stub := shim.NewMockStub("test", tcc)

	// Init A=345 B=456
	checkInit(t, stub, [][]byte{[]byte("init"), []byte("Alice"), []byte("100"), []byte("Bob"), []byte("200")})

	// Query A
	checkQuery(t, stub, [][]byte{[]byte("invoke"), []byte("query"), []byte("Alice")}, "100")

	// Query B
	checkQuery(t, stub, [][]byte{[]byte("invoke"), []byte("query"), []byte("Bob")}, "200")
}

//Test if invoke works as intended
func TestChaincode_Invoke(t *testing.T) {
	tcc := new(main.SimpleChaincode)
	stub := shim.NewMockStub("test", tcc);

	// Init Alice=100 Bob=200
	checkInit(t, stub, [][]byte{[]byte("init"), []byte("Alice"), []byte("100"), []byte("Bob"), []byte("200")})


	// Invoke Alice -> Bob for 50
	checkInvoke(t, stub, [][]byte{[]byte("invoke"), []byte("pay"), []byte("Alice"), []byte("Bob"), []byte("50")})
	checkQuery(t, stub, [][]byte{[]byte("invoke"), []byte("query"), []byte("Alice")}, "50")
	checkQuery(t, stub, [][]byte{[]byte("invoke"), []byte("query"), []byte("Bob")}, "250")

	checkInvoke(t, stub, [][]byte{[]byte("invoke"), []byte("pay"), []byte("Bob"), []byte("Alice"), []byte("125")})
	checkQuery(t, stub, [][]byte{[]byte("invoke"), []byte("query"), []byte("Alice")}, "175")
	checkQuery(t, stub, [][]byte{[]byte("invoke"), []byte("query"), []byte("Bob")}, "125")
}

