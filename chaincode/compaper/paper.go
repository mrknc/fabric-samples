package main

import (
	"encoding/json"
	"fmt"
	"os"

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
	Maturity     string `json:"maturity`
	FaceValue    string `json:"faceValue"` // TODO this doesn't seem right
	CurrentState string `json:"currentState"`
}

func (p Paper) compositeKey() string {
	return "papernet" + p.Issuer + p.Paper
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
		return shim.Error("Incorrect number of arguments. Expecting 5")
	}

	// TODO they appear not to do any checks whatsoever: what if there is a paper already issued?
	// assigned values should make sense, like dates, face value not being negative
	paper := Paper{Issuer: args[0], Paper: args[1], Owner: args[0],
		Issue: args[2], Maturity: args[3], FaceValue: args[4], CurrentState: isISSUED}

	paperAsBytes, _ := json.Marshal(paper)
	fmt.Println(paper)
	APIstub.PutState(paper.compositeKey(), paperAsBytes)

	return shim.Success(nil)
}

func (s *SmartContract) redeemPaper(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	return shim.Success(nil)
}

func (s *SmartContract) buyPaper(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	return shim.Success(nil)
}

func main() {
	x := Paper{"gogi", "00001", "gogi", "2020-01-01", "2020-07-01", "100", "ISSUED"}
	fmt.Println(x)
	j, err := json.Marshal(x)
	fmt.Println(j, err)
	os.Stdout.Write(j)
	y := Paper{}
	json.Unmarshal(j, &y)
	fmt.Println("\nyo", y)
	args := []string{"mogi", "00002", "mogi", "2020-01-01", "2020-09-01", "200", "ISSUED"}
	fmt.Println(args)
	z := Paper{Issuer: args[0], Paper: args[1]}
	fmt.Println(z)
	fmt.Println(z.compositeKey())

}
