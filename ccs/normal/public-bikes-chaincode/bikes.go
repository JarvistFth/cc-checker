package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"regexp"
	"strconv"
	"strings"
)

var logger = shim.NewLogger("GiouChaincode")

// Participant types - Each participant type is mapped to an integer which we use to compare to the value stored in a
// user's eCert.
// CURRENT WORKAROUND USES ROLES CHANGE WHEN OWN USERS CAN BE CREATED SO THAT IT READ 1, 2, 3, 4
const AUTHORITY = "regulator"
const TOWNSHIP = "township"
const PRIVATE_ENTITY = "private"
const STATION = "station"

// Status types - Asset lifecycle is broken down into 4 statuses, this is part of the business logic to determine what can
// be done to the bike at points in it's lifecycle.
const STATE_TEMPLATE = 0
const STATE_TOWNSHIP = 1
const STATE_PRIVATE_OWNERSHIP = 2
const STATE_STATION = 3

// SimpleChaincode is a blank struct for use with Shim (A HyperLedger included go file used for get/put state
// and other HyperLedger functions).
type SimpleChaincode struct {
}

// Bike defines the structure for a bike object. JSON on right tells it what JSON fields to map to
// that element when reading a JSON object into the struct e.g. JSON make -> Struct Make.
type Bike struct {
	Model  string `json:"model"`
	Reg    string `json:"reg"`   //registration
	BTWSH  int    `json:"BTWSH"` //"BikeTownship" : An integer value given by the Township that the bike belongs
	Owner  string `json:"owner"`
	Status int    `json:"status"` //values : 0(Regulator), 1(Township), 2(Private Entity), 3(Station)
	Colour string `json:"colour"`
	crtID  string `json:"crtID"` //"createID" : unique key
}

// crtHolder defines the structure that holds all the crtIDs for bikes that have been created.
// Used as an index when querying all bikes.
type crtHolder struct {
	crts []string `json:"crts"`
}

//	UserAndECert is the struct for storing the JSON of a user and their ecert.
type UserAndECert struct {
	Identity string `json:"identity"`
	eCert    string `json:"ecert"`
}

//	Init function is called when the user deploys the chaincode.
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	_, args := stub.GetFunctionAndParameters()
	var crtIDs crtHolder

	bytes, err := json.Marshal(crtIDs)
	if err != nil {
		return shim.Error("Error creating crtHolder record")
	}

	err = stub.PutState("crtIDs", bytes)
	if err != nil {
		return shim.Error("Error storing bike to blockchain.")
	}

	for i := 0; i < len(args); i = i + 2 {
		t.addEcert(stub, args[i], args[i+1])
	}

	return shim.Success([]byte("ok"))
}

// getEcert takes the name passed and calls out to the REST API for HyperLedger to retrieve the ecert
// for that user. Returns the ecert as retrived including html encoding.
func (t *SimpleChaincode) getEcert(stub shim.ChaincodeStubInterface, name string) pb.Response {

	ecert, err := stub.GetState(name)
	if err != nil {
		return shim.Error("Couldn't retrieve ecert for user " + name)
	}

	return shim.Success(ecert)
}

// addEcert - Adds a new ecert and user pair to the table of ecerts
func (t *SimpleChaincode) addEcert(stub shim.ChaincodeStubInterface, name string, ecert string) pb.Response {

	err := stub.PutState(name, []byte(ecert))
	if err != nil {
		return shim.Error("Error storing eCert for user " + name + " identity: " + ecert)
	}

	return shim.Success(nil)
}

// getUsername retrieves the username of the user who invoked the chaincode and returns the username as a string.
func (t *SimpleChaincode) getUsername(stub shim.ChaincodeStubInterface) (string, error) {

	username, err := stub.GetState("username")
	if err != nil {
		return "", errors.New("Couldn't get attribute 'username'. Error: " + err.Error())
	}

	return string(username), nil
}

// checkAffiliation takes an ecert as a string, decodes it to remove html encoding then parses it and checks the
// certificates common name. The affiliation is stored as part of the common name.
func (t *SimpleChaincode) checkAffiliation(stub shim.ChaincodeStubInterface) (string, error) {

	affiliation, err := stub.GetState("role")
	if err != nil {
		return "", errors.New("Couldn't get attribute 'role'. Error: " + err.Error())
	}

	return string(affiliation), nil
}

// getCallerData calls the getEcert and checkRole functions and returns the ecert and role for the
// name passed.
func (t *SimpleChaincode) getCallerData(stub shim.ChaincodeStubInterface) (string, string, error) {
	user, err := t.getUsername(stub)

	affiliation, err := t.checkAffiliation(stub)
	if err != nil {
		return "", "", err
	}

	return user, affiliation, nil
}

// retrieveCrt gets the state of the data at crtID in the ledger then converts it from the stored
// JSON into the Bike struct for use in the contract. Returns the Bike struct. Returns empty v if it errors.
func (t *SimpleChaincode) retrieveCrt(stub shim.ChaincodeStubInterface, crtID string) (Bike, error) {
	var v Bike

	bytes, err := stub.GetState(crtID)
	if err != nil {
		fmt.Printf("retrieveCrt: Failed to invoke bikeCode: %s", err)
		return v, errors.New("retrieveCrt: Error retrieving bike with crtID = " + crtID)
	}

	err = json.Unmarshal(bytes, &v)
	if err != nil {
		fmt.Printf("retrieveCrt: Corrupt bike record "+string(bytes)+": %s", err)
		return v, errors.New("retrieveCrt: Corrupt bike record" + string(bytes))
	}

	return v, nil
}

// saveChanges writes to the ledger the Bike struct passed in a JSON format. Uses the shim file's
// method 'PutState'.
func (t *SimpleChaincode) saveChanges(stub shim.ChaincodeStubInterface, v Bike) (bool, error) {
	bytes, err := json.Marshal(v)
	if err != nil {
		fmt.Printf("saveChanges: Error converting bike record: %s", err)
		return false, errors.New("Error converting bike record")
	}

	err = stub.PutState(v.crtID, bytes)
	if err != nil {
		fmt.Printf("saveChanges: Error storing bike record: %s", err)
		return false, errors.New("Error storing bike record")
	}

	return true, nil
}

// Invoke is called on chaincode invoke. Takes a function name passed and calls that function. Converts some
// initial arguments passed to other things for use in the called function e.g. name -> ecert
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	caller, callerAffiliation, err := t.getCallerData(stub)
	if err != nil {
		return shim.Error("Error retrieving caller information")
	}

	if function == "createBike" {
		return t.createBike(stub, caller, callerAffiliation, args[0])
	} else { // If the function is not a create then there must be a bike so we need to retrieve the bike.
		argPos := 1

		v, err := t.retrieveCrt(stub, args[argPos])
		if err != nil {
			fmt.Printf("invoke: Error retrieving crt: %s", err)
			return shim.Error("Error retrieving crt")
		}

		if strings.Contains(function, "update") == false {
			if function == "authorityToTownship" {
				return t.authorityToTownship(stub, v, caller, callerAffiliation, args[0], "township")
			} else if function == "townshipToStation" {
				return t.townshipToStation(stub, v, caller, callerAffiliation, args[0], "station")
			} else if function == "stationToPrivate" {
				return t.stationToPrivate(stub, v, caller, callerAffiliation, args[0], "private")
			} else if function == "privateToStation" {
				return t.privateToStation(stub, v, caller, callerAffiliation, args[0], "station")
			}
		} else if function == "updateModel" {
			return t.updateModel(stub, v, caller, callerAffiliation, args[0])
		} else if function == "updateRegistration" {
			return t.updateRegistration(stub, v, caller, callerAffiliation, args[0])
		} else if function == "updateBtwsh" {
			return t.updateBtwsh(stub, v, caller, callerAffiliation, args[0])
		} else if function == "updateColour" {
			return t.updateColour(stub, v, caller, callerAffiliation, args[0])
		}

		return shim.Error("Function of the name " + function + " doesn't exist.")

	}
}

// Query is called on chaincode query. Takes a function name passed and calls that function. Passes the
// initial arguments passed are passed on to the called function.
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) pb.Response {
	caller, callerAffiliation, err := t.getCallerData(stub)
	if err != nil {
		fmt.Printf("query: Error retrieving caller details", err)
		return shim.Error("query: Error retrieving caller details: " + err.Error())
	}

	if function == "getBikeDetails" {
		if len(args) != 1 {
			fmt.Printf("Incorrect number of arguments passed")
			return shim.Error("query: Incorrect number of arguments passed")
		}

		v, err := t.retrieveCrt(stub, args[0])

		if err != nil {
			fmt.Printf("query: Error retrieving crt: %s", err)
			return shim.Error("query: Error retrieving crt " + err.Error())
		}

		return t.getBikeDetails(stub, v, caller, callerAffiliation)
	} else if function == "checkUniqueCrt" {
		return t.checkUniqueCrt(stub, args[0], caller, callerAffiliation)
	} else if function == "getBikes" {
		return t.getBikes(stub, caller, callerAffiliation)
	} else if function == "getEcert" {
		return t.getEcert(stub, args[0])
	} else if function == "ping" {
		return shim.Success([]byte("hello"))
	}

	return shim.Error("Received unknown function invocation")
}

// createBike creates the initial JSON for the bike and then saves it to the ledger.
func (t *SimpleChaincode) createBike(stub shim.ChaincodeStubInterface, caller string, callerAffiliation string,
	crtID string) pb.Response {
	var v Bike

	crt_ID := "\"crtID\":\"" + crtID + "\", " // Variables to define the JSON
	btwsh := "\"BTWSH\":0, "
	model := "\"Model\":\"UNDEFINED\", "
	reg := "\"Reg\":\"UNDEFINED\", "
	owner := "\"Owner\":\"" + caller + "\", "
	colour := "\"Colour\":\"UNDEFINED\", "
	status := "\"Status\":0 "

	bike_json := "{" + crt_ID + btwsh + model + reg + owner + colour + status + "}" // Concatenates the variables to create the total JSON object

	matched, err := regexp.Match("^[A-z][A-z][0-9]{7}", []byte(crtID)) // matched = true if the crtID passed fits format of two letters followed by seven digits
	if err != nil {
		fmt.Printf("createBike: Invalid crtID: %s", err)
		return shim.Error("Invalid crtID")
	}

	if crt_ID == "" || matched == false {
		fmt.Printf("createBike: Invalid crtID provided")
		return shim.Error("Invalid crtID provided")
	}

	err = json.Unmarshal([]byte(bike_json), &v) // Convert the JSON defined above into a bike object for go
	if err != nil {
		return shim.Error("Invalid JSON object")
	}

	record, err := stub.GetState(v.crtID) // If not an error then a record exists so cant create a new bike with this crtID as it must be unique
	if record != nil {
		return shim.Error("Bike already exists")
	}

	if callerAffiliation != AUTHORITY { // Only the regulator can create a new crt
		return shim.Error(fmt.Sprintf("Permission Denied. createBike. %v === %v", callerAffiliation, AUTHORITY))
	}

	_, err = t.saveChanges(stub, v)
	if err != nil {
		fmt.Printf("createBike: Error saving changes: %s", err)
		return shim.Error("Error saving changes")
	}

	bytes, err := stub.GetState("crtIDs")
	if err != nil {
		return shim.Error("Unable to get crtIDs")
	}

	var crtIDs crtHolder
	err = json.Unmarshal(bytes, &crtIDs)
	if err != nil {
		return shim.Error("Corrupt crtHolder record")
	}

	crtIDs.crts = append(crtIDs.crts, crtID)

	bytes, err = json.Marshal(crtIDs)
	if err != nil {
		fmt.Print("Error creating crtHolder record")
	}

	err = stub.PutState("crtIDs", bytes)
	if err != nil {
		return shim.Error("Unable to put the state")
	}

	return shim.Success(nil)

}

// authorityToTownship transfer ownership from authority to township.
func (t *SimpleChaincode) authorityToTownship(stub shim.ChaincodeStubInterface, v Bike, caller string,
	callerAffiliation string, recipientName string, recipientAffiliation string) pb.Response {

	if v.Status == STATE_TEMPLATE && v.Owner == caller && callerAffiliation == AUTHORITY &&
		recipientAffiliation == TOWNSHIP { // If the roles and users are ok

		v.Owner = recipientName   // then make the owner the new owner
		v.Status = STATE_TOWNSHIP // and mark it in the state of township

	} else { // Otherwise if there is an error
		fmt.Printf("authorityToTownship: Permission Denied")
		return shim.Error(fmt.Sprintf("Permission Denied. authorityToTownship. %v %v === %v,"+
			" %v === %v, %v === %v, %v === %v", v, v.Status, STATE_TEMPLATE, v.Owner, caller, callerAffiliation,
			AUTHORITY, recipientAffiliation, TOWNSHIP))
	}

	_, err := t.saveChanges(stub, v) // Write new state
	if err != nil {
		fmt.Printf("authorityToTownship: Error saving changes: %s", err)
		return shim.Error("Error saving changes")
	}

	return shim.Success(nil)
}

// townshipToStation transfer ownership from township to station.
func (t *SimpleChaincode) townshipToStation(stub shim.ChaincodeStubInterface, v Bike, caller string,
	callerAffiliation string, recipientName string, recipientAffiliation string) pb.Response {

	if v.Model == "UNDEFINED" || v.Reg == "UNDEFINED" || v.Colour == "UNDEFINED" || v.BTWSH == 0 { //If any part of the bike is undefined it has not been fully defined so cannot be sent
		fmt.Printf("townshipToStation: Bike not fully defined")
		return shim.Error(fmt.Sprintf("Bike not fully defined. %v", v))
	}

	if v.Status == STATE_TOWNSHIP && v.Owner == caller && callerAffiliation == TOWNSHIP && recipientAffiliation == STATION {

		v.Owner = recipientName
		v.Status = STATE_STATION
	} else {
		return shim.Error(fmt.Sprintf("Permission Denied. townshipToStation. %v %v === %v, %v === %v,"+
			" %v === %v, %v === %v", v, v.Status, STATE_TOWNSHIP, v.Owner, caller, callerAffiliation, TOWNSHIP,
			recipientAffiliation, STATION))
	}

	_, err := t.saveChanges(stub, v)
	if err != nil {
		fmt.Printf("townshipToStation: Error saving changes: %s", err)
		return shim.Error("Error saving changes")
	}

	return shim.Success(nil)

}

// stationToPrivate transfer ownership from station to private entity.
func (t *SimpleChaincode) stationToPrivate(stub shim.ChaincodeStubInterface, v Bike, caller string,
	callerAffiliation string, recipientName string, recipientAffiliation string) pb.Response {

	if v.Status == STATE_STATION && v.Owner == caller && callerAffiliation == STATION &&
		recipientAffiliation == PRIVATE_ENTITY {

		v.Owner = recipientName
	} else {
		return shim.Error(fmt.Sprintf("Permission Denied. stationToPrivate. %v %v === %v, %v === %v,"+
			" %v === %v, %v === %v", v, v.Status, STATE_STATION, v.Owner, caller, callerAffiliation, STATION,
			recipientAffiliation, PRIVATE_ENTITY))
	}

	_, err := t.saveChanges(stub, v)
	if err != nil {
		fmt.Printf("stationToPrivate: Error saving changes: %s", err)
		return shim.Error("Error saving changes")
	}

	return shim.Success(nil)
}

// privateToStation transfer ownership from private entity to station.
func (t *SimpleChaincode) privateToStation(stub shim.ChaincodeStubInterface, v Bike, caller string,
	callerAffiliation string, recipientName string, recipientAffiliation string) pb.Response {

	if v.Status == STATE_PRIVATE_OWNERSHIP && v.Owner == caller && callerAffiliation == PRIVATE_ENTITY &&
		recipientAffiliation == STATION {

		v.Owner = recipientName
	} else {
		return shim.Error(fmt.Sprintf("Permission denied. privateToStation."+
			" %v === %v, %v === %v, %v === %v, %v === %v", v.Status, STATE_PRIVATE_OWNERSHIP, v.Owner, caller,
			callerAffiliation, PRIVATE_ENTITY, recipientAffiliation, STATION))
	}

	_, err := t.saveChanges(stub, v)
	if err != nil {
		fmt.Printf("privateToStation: Error saving changes: %s", err)
		return shim.Error("Error saving changes")
	}

	return shim.Success(nil)

}

// updateBtwsh updates bikes township entity.
func (t *SimpleChaincode) updateBtwsh(stub shim.ChaincodeStubInterface, v Bike, caller string,
	callerAffiliation string, new_value string) pb.Response {

	new_btwsh, err := strconv.Atoi(string(new_value)) // will return an error if the new btwsh contains non numerical chars
	if err != nil || len(string(new_value)) != 15 {
		return shim.Error("Invalid value passed for new BTWSH")
	}
	// Can't change the BTWSH after its initial assignment (BTWSH==0)
	if v.Status == STATE_TOWNSHIP && v.Owner == caller && callerAffiliation == TOWNSHIP && v.BTWSH == 0 {

		v.BTWSH = new_btwsh // Update to the new value
	} else {
		return shim.Error(fmt.Sprintf("Permission denied. updateBtwsh %v %v %v %v %v", v.Status,
			STATE_TOWNSHIP, v.Owner, caller, v.BTWSH))
	}

	_, err = t.saveChanges(stub, v) // Save the changes in the blockchain
	if err != nil {
		fmt.Printf("updateBtwsh: Error saving changes: %s", err)
		return shim.Error("Error saving changes")
	}

	return shim.Success(nil)

}

// updateRegistration updates registration.
func (t *SimpleChaincode) updateRegistration(stub shim.ChaincodeStubInterface, v Bike, caller string,
	callerAffiliation string, new_value string) pb.Response {

	if v.Owner == caller {
		v.Reg = new_value
	} else {
		return shim.Error(fmt.Sprint("Permission denied. updateRegistration"))
	}

	_, err := t.saveChanges(stub, v)
	if err != nil {
		fmt.Printf("updateRegistration: Error saving changes: %s", err)
		return shim.Error("Error saving changes")
	}

	return shim.Success(nil)
}

// updateColour updates colour of bike asset.
func (t *SimpleChaincode) updateColour(stub shim.ChaincodeStubInterface, v Bike, caller string,
	callerAffiliation string, new_value string) pb.Response {

	if v.Owner == caller && callerAffiliation == TOWNSHIP {
		v.Colour = new_value
	} else {
		return shim.Error(fmt.Sprint("Permission denied. updateColour %t %t"+v.Owner == caller, callerAffiliation == TOWNSHIP))
	}

	_, err := t.saveChanges(stub, v)
	if err != nil {
		fmt.Printf("updateColour: Error saving changes: %s", err)
		return shim.Error("Error saving changes")
	}

	return shim.Success(nil)
}

// updateModel updates model of bike asset.
func (t *SimpleChaincode) updateModel(stub shim.ChaincodeStubInterface, v Bike, caller string,
	callerAffiliation string, new_value string) pb.Response {

	if v.Status == STATE_TOWNSHIP && v.Owner == caller && callerAffiliation == TOWNSHIP {
		v.Model = new_value
	} else {
		return shim.Error(fmt.Sprint("Permission denied. updateModel %t %t"+v.Owner == caller, callerAffiliation == TOWNSHIP))

	}

	_, err := t.saveChanges(stub, v)
	if err != nil {
		fmt.Printf("updateModel: Error saving changes: %s", err)
		return shim.Error("Error saving changes")
	}

	return shim.Success(nil)
}

// getBikeDetails gets bike details.
func (t *SimpleChaincode) getBikeDetails(stub shim.ChaincodeStubInterface, v Bike, caller string,
	callerAffiliation string) pb.Response {

	bytes, err := json.Marshal(v)
	if err != nil {
		return shim.Error("getBikeDetails: Invalid bike object")
	}

	if v.Owner == caller ||
		callerAffiliation == AUTHORITY {
		return shim.Success(bytes)
	} else {
		return shim.Error("Permission Denied. getBikeDetails")
	}

}

// getBikes gets bikes.
func (t *SimpleChaincode) getBikes(stub shim.ChaincodeStubInterface, caller string, callerAffiliation string) pb.Response {
	bytes, err := stub.GetState("crtIDs")
	if err != nil {
		return shim.Error("Unable to get crtIDs")
	}

	var crtIDs crtHolder
	err = json.Unmarshal(bytes, &crtIDs)
	if err != nil {
		return shim.Error("Corrupt crtHolder")
	}

	result := "["
	var v Bike

	for _, crt := range crtIDs.crts {

		v, err = t.retrieveCrt(stub, crt)
		if err != nil {
			return shim.Error("Failed to retrieve crt")
		}

		temp := t.getBikeDetails(stub, v, caller, callerAffiliation)
		if temp.Status == shim.OK {
			result += string(temp.Payload) + ","
		}
	}

	if len(result) == 1 {
		result = "[]"
	} else {
		result = result[:len(result)-1] + "]"
	}

	return shim.Success([]byte(result))
}

// checkUniqueCrt checks unique crt.
func (t *SimpleChaincode) checkUniqueCrt(stub shim.ChaincodeStubInterface, crt string, caller string,
	callerAffiliation string) pb.Response {

	_, err := t.retrieveCrt(stub, crt)
	if err == nil {
		return shim.Success([]byte("false\n crt is not unique"))
	} else {
		return shim.Success([]byte("true"))
	}
}

// main starts up the chaincode
func main() {

	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Chaincode: %s", err)
	}
}
