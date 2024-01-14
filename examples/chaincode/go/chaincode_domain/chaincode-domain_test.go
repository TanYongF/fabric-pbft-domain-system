package main

import (
	"fmt"
	"testing"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

func checkInit(t *testing.T, stub *shim.MockStub, args []string) {
	_, err := stub.MockInit("1", "init", args)
	if err != nil {
		fmt.Println("Init failed", err)
		t.FailNow()
	}
}

func checkState(t *testing.T, stub *shim.MockStub, name string, value string) {
	bytes := stub.State[name]
	if bytes == nil {
		fmt.Println("State", name, "failed to get value")
		t.FailNow()
	}
	if string(bytes) != value {
		fmt.Println("State value", name, "was not", value, "as expected")
		t.FailNow()
	}
}

func checkQuery(t *testing.T, stub *shim.MockStub, function string, name string, value string) {
	bytes, err := stub.MockQuery(function, []string{name})
	if err != nil {
		fmt.Println("Query", name, "failed", err)
		t.FailNow()
	}
	if bytes == nil {
		fmt.Println("Query", name, "failed to get value")
		if value != "" {
			t.FailNow()
		}
	}
	if string(bytes) != value {
		fmt.Println("Query value", name, "was not", value, "as expected")
		t.FailNow()
	}
}

func checkInvoke(t *testing.T, stub *shim.MockStub, args []string) {
	_, err := stub.MockInvoke("1", "query", args)
	if err != nil {
		fmt.Println("Invoke", args, "failed", err)
		t.FailNow()
	}
}

func TestExample02_Init(t *testing.T) {
	scc := new(SimpleChaincode)
	stub := shim.NewMockStub("ex10", scc)

	// Init A="" B=""
	checkInit(t, stub, []string{"A", "", "B", ""})

	checkState(t, stub, "A", "{}")
	checkState(t, stub, "B", "{}")
}

func TestExample02_Query(t *testing.T) {
	scc := new(SimpleChaincode)
	stub := shim.NewMockStub("ex10", scc)

	// Init A have no domain, B have no domain
	checkInit(t, stub, []string{"A", "", "B", ""})
	// Query A
	checkQuery(t, stub, "getDomainsByOwner", "A", "{}")
	// Query B
	checkQuery(t, stub, "getOwnerByDomain", "google.com", "")
}

func TestExample02_Invoke(t *testing.T) {
	scc := new(SimpleChaincode)
	stub := shim.NewMockStub("ex10", scc)

	// A and B have no domain
	checkInit(t, stub, []string{"A", "", "B", ""})

	// A get domain google.com
	checkInvoke(t, stub, []string{"A", "google.com"})
	checkInvoke(t, stub, []string{"A", "hello.com"})
	checkInvoke(t, stub, []string{"B", "baidu.com"})
	checkQuery(t, stub, "getDomainsByOwner", "A", "{\"google.com\":{},\"hello.com\":{}}")
	checkQuery(t, stub, "getOwnerByDomain", "google.com", "A")

}
