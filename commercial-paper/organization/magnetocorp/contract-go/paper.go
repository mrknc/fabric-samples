package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

const (
	isISSUED   = "ISSUED"
	isTRADING  = "TRADING"
	isREDEEMED = "REDEEMED"
)

type Paper struct {
	Issuer       string `json:"issuer"`
	Paper        string `json:"paper"`
	Owner        string `json:"owner"`
	Issue        string `json:"time"` // TODO this doesn't seem right
	Maturity     string `json:"maturity"`
	FaceValue    string `json:"faceValue"` // TODO this doesn't seem right
	CurrentState string `json:"currentState"`
}

func compositeKey(issuer, paper string) string {
	return "papernet" + issuer + paper
}

// Define the Smart Contract structure
type SmartContract struct {
}

/*
* Called at instantiation
 */
func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}

/*
* The Invoke method is called as a result of an application request to run the Smart Contract
* The calling application program has also specified the particular smart contract function to be called, with arguments
 */
func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {
	// Retrieve the requested Smart Contract function and arguments
	function, args := APIstub.GetFunctionAndParameters()
	// Route to the appropriate handler function to interact with the ledger appropriately
	if function == "issuePaper" {
		return s.issuePaper(APIstub, args)
	} else if function == "redeemPaper" {
		return s.redeemPaper(APIstub, args)
	} else if function == "buyPaper" {
		return s.buyPaper(APIstub, args)
	}
	return shim.Error("Invalid Smart Contract function name.")
}

func (s *SmartContract) issuePaper(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) != 5 {
		return shim.Error("Incorrect number of arguments. Expecting 5: issuer, paper, issue date, maturity date, face value")
	}

	// TODO they appear not to do any checks whatsoever: what if there is a paper already issued?
	// assigned values should make sense, like dates, face value not being negative
	p := Paper{Issuer: args[0], Paper: args[1], Owner: args[0],
		Issue: args[2], Maturity: args[3], FaceValue: args[4], CurrentState: isISSUED}
	paperAsBytes, _ := json.Marshal(p)
	APIstub.PutState(compositeKey(p.Issuer, p.Paper), paperAsBytes)

	return shim.Success(nil)
}

func (s *SmartContract) redeemPaper(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3: issuer, paper, redeeming owner")
	}
	key := compositeKey(args[0], args[1])
	paperAsBytes, _ := APIstub.GetState(key) //TODO err handling
	p := Paper{}
	json.Unmarshal(paperAsBytes, &p)
	if p.Owner != args[2] {
		return shim.Error("Expected owner: " + args[2] + " actual: " + p.Owner)
	}
	if p.CurrentState == isREDEEMED {
		return shim.Error("Paper " + key + " was already redeemed.")
	}
	p.Owner = p.Issuer
	p.CurrentState = isREDEEMED
	paperAsBytes, _ = json.Marshal(p)
	APIstub.PutState(key, paperAsBytes)
	return shim.Success(nil)
}

func (s *SmartContract) buyPaper(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments. Expecting 4: issuer, paper, owner, new owner")
	}
	key := compositeKey(args[0], args[1])
	paperAsBytes, _ := APIstub.GetState(key) //TODO err handling
	p := Paper{}
	json.Unmarshal(paperAsBytes, &p)
	if p.Owner != args[2] {
		return shim.Error("Expected owner: " + args[2] + " actual: " + p.Owner)
	}
	p.Owner = args[3]
	p.CurrentState = isTRADING
	paperAsBytes, _ = json.Marshal(p)
	APIstub.PutState(key, paperAsBytes)
	return shim.Success(nil)
}

func main() {
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
