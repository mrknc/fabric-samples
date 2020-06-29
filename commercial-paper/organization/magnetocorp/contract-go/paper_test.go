package main

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

func checkInit(t *testing.T, stub *shim.MockStub, args [][]byte) {
	res := stub.MockInit("mockstub", args)
	if res.Status != shim.OK {
		fmt.Println("Init failed", string(res.Message))
		t.FailNow()
	}
}

func checkState(t *testing.T, stub *shim.MockStub, name string, paper Paper) {
	bytes := stub.State[name]
	if bytes == nil {
		fmt.Println("State", name, "failed to get value")
		t.FailNow()
	}
	if pbytes, _ := json.Marshal(paper); string(bytes) != string(pbytes) {
		fmt.Println("State value", name, "was different from", paper)
		fmt.Println(string(bytes))
		t.FailNow()
	}
}

func checkInvoke(t *testing.T, stub *shim.MockStub, args [][]byte) {
	res := stub.MockInvoke("1", args)
	if res.Status != shim.OK {
		fmt.Println("Invoke", args, "failed", string(res.Message))
		t.FailNow()
	}
}

func TestExample02_Init(t *testing.T) {
	scc := new(SmartContract)
	stub := shim.NewMockStub("cpgo", scc)

	checkInit(t, stub, nil)
}

func TestExample02_Invoke(t *testing.T) {
	scc := new(SmartContract)
	stub := shim.NewMockStub("cpgo", scc)

	checkInit(t, stub, nil)

	checkInvoke(t, stub, [][]byte{[]byte("issuePaper"), []byte("GOGI"), []byte("00001"), []byte("2020-01-01"),
		[]byte("2020-07-01"), []byte("1000")})
	fmt.Println(stub.State)
	tp := Paper{"GOGI", "00001", "GOGI", "2020-01-01", "2020-07-01", "1000", "ISSUED"}
	checkState(t, stub, compositeKey("GOGI", "00001"), tp)

	checkInvoke(t, stub, [][]byte{[]byte("buyPaper"), []byte("GOGI"), []byte("00001"), []byte("GOGI"),
		[]byte("MOGI")})
	tp.Owner = "MOGI"
	tp.CurrentState = "TRADING"
	checkState(t, stub, compositeKey("GOGI", "00001"), tp)

	checkInvoke(t, stub, [][]byte{[]byte("redeemPaper"), []byte("GOGI"), []byte("00001"), []byte("MOGI")})
	tp.Owner = "GOGI"
	tp.CurrentState = isREDEEMED
	checkState(t, stub, compositeKey("GOGI", "00001"), tp)

}
