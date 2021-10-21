package randomcc

import (
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
	"math/rand"
)

// SimpleAsset implements a simple chaincode to manage an asset
type SimpleAsset struct {

}

// Init is called during chaincode instantiation to initialize any
// data. Note that chaincode upgrade also calls this function to reset
// or to migrate data.
func (t *SimpleAsset) Init(stub shim.ChaincodeStubInterface) peer.Response {
	// Get the args from the transaction proposal
	_,args := stub.GetFunctionAndParameters()
	if len(args) != 2 {
		return shim.Error(fmt.Sprintf("incorrect args, len(args):%d", len(args)))
	}

	// Set up any variables or assets here by calling stub.PutState()

	// We store the key and the value on the ledger
	err := stub.PutState(args[0], []byte(args[1]))
	if err != nil {
		return shim.Error(fmt.Sprintf("Failed to create asset: %s", args[0]))
	}
	return shim.Success([]byte(fmt.Sprintf("create key: %s, success",args[0])))
}

// Invoke is called per transaction on the chaincode. Each transaction is
// either a 'get' or a 'set' on the asset created by Init function. The Set
// method may create a new asset by specifying a new key-value pair.
func (t *SimpleAsset) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	// Extract the function and args from the transaction proposal

	var result string
	var err error
	r := getRandom()

	result,err = set(stub,"randkey",[]byte(fmt.Sprintf("%d",r)))
	if err != nil {
		return shim.Error(err.Error())
	}

	// Return the result as success payload
	return shim.Success([]byte(result))
}

func getRandom() int {
	return rand.Int()
}

// Set stores the asset (both key and value) on the ledger. If the key exists,
// it will override the value with the new one
func set(stub shim.ChaincodeStubInterface, key string, value []byte) (string, error) {


	err := stub.PutState(key, value)
	if err != nil {
		return "", fmt.Errorf("Failed to set asset: %s", key)
	}
	return string(value), nil
}

// Get returns the value of the specified asset key
func get(stub shim.ChaincodeStubInterface, key string) (string, error) {
	if len(key) == 0 {
		return "", fmt.Errorf("Incorrect arguments. Expecting a key")
	}

	value, err := stub.GetState(key)
	if err != nil {
		return "", fmt.Errorf("Failed to get asset: %s with error: %s", key, err)
	}
	if value == nil {
		return "", fmt.Errorf("Asset not found: %s", key)
	}
	return string(value), nil
}

// main function starts up the chaincode in the container during instantiate
func main() {
	if err := shim.Start(new(SimpleAsset)); err != nil {
		fmt.Printf("Error starting SimpleAsset chaincode: %s", err)
	}
}