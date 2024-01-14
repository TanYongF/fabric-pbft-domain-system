package main

//WARNING - this chaincode's ID is hard-coded in chaincode_example04 to illustrate one way of
//calling chaincode from a chaincode. If this example is modified, chaincode_example04.go has
//to be modified as well with the new ID of chaincode_example02.
//chaincode_example05 show's how chaincode ID can be passed in as a parameter instead of
//hard-coding.

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}
type void struct{}
type set map[string]void

// strategies 用于存储查询策略
var strategies map[string]func(stub shim.ChaincodeStubInterface, args []string) ([]byte, error)

const QUERY_OWNER_BY_DOMAIN = "getOwnerByDomain"
const QUERY_DOMAINS_BY_OWNER = "getDomainsByOwner"

func init() {
	strategies = make(map[string]func(stub shim.ChaincodeStubInterface, args []string) ([]byte, error))
	strategies[QUERY_DOMAINS_BY_OWNER] = GetOwnerByDomain
	strategies[QUERY_OWNER_BY_DOMAIN] = GetDomainsByOwner
}

func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	var A, B string    // Entities
	var Aval, Bval set // Asset holdings
	var err error

	if len(args) != 4 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2")
	}

	// Initialize the chaincode
	A = args[0]
	Aval = make(set)
	B = args[2]
	Bval = make(set)

	// map2json
	AvalStr, _ := json.Marshal(Aval)
	BvalStr, _ := json.Marshal(Bval)

	// Write the state to the ledger
	err = stub.PutState(A, AvalStr)
	//err = stub.PutState()
	if err != nil {
		fmt.Printf("put state error")
		return nil, err
	}

	err = stub.PutState(B, BvalStr)
	if err != nil {
		fmt.Printf("put state error")
		return nil, err
	}

	return nil, nil
}

// Transaction makes payment of X units from A to B
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	if function == "delete" {
		// Deletes an entity from its state
		return t.delete(stub, args)
	}

	var A string             // owner, args[0]
	var requestDomain string // domain, args[1]
	var err error

	if len(args) != 2 {
		return nil, errors.New("incorrect number of arguments. Expecting 2")
	}

	A = args[0]
	requestDomain = args[1]

	// 1. 判断这个域名有无拥有者

	// 2. 设置域名-拥有者的关系
	err = stub.PutState(requestDomain, []byte(A))
	if err != nil {
		return nil, errors.New("failed to update the domain-owner relation")
	}

	// 3. 更新拥有者-域名集合的关系
	domainList, err := stub.GetState(A)
	if err != nil {
		return nil, errors.New("failed to get domain list")
	}
	var domains set
	err = json.Unmarshal(domainList, &domains)
	if err != nil {
		return nil, errors.New("failed to unmarshal the domain list")
	}
	var member void
	domains[requestDomain] = member
	domainsJson, err := json.Marshal(domains)
	if err != nil {
		return nil, errors.New("failed to marshal the domains")
	}

	err = stub.PutState(A, domainsJson)
	if err != nil {
		return nil, errors.New("failed to put state")
	}

	return nil, nil
}

// Deletes an entity from state
func (t *SimpleChaincode) delete(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != 1 {
		return nil, errors.New("incorrect number of arguments. Expecting 1")
	}

	A := args[0]

	// Delete the key from the state in ledger
	err := stub.DelState(A)
	if err != nil {
		return nil, errors.New("Failed to delete state")
	}

	return nil, nil
}

// Query callback representing the query of a chaincode
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	if len(function) == 0 || strategies[function] == nil {
		return nil, errors.New("invalid query function name. Expecting \"query\"")
	}
	if len(args) != 1 {
		return nil, errors.New("incorrect number of arguments, Expecting 1 args")
	}

	if response, err := strategies[function](stub, args); err != nil {
		return nil, err
	} else {
		fmt.Printf("Query Response:%s\n", response)
		return response, nil
	}

}

// GetOwnerByDomain 通过域名获取用户名
func GetOwnerByDomain(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	domain := args[0]
	owner, err := stub.GetState(domain)
	if err != nil {
		jsonResp := "{\"Error\":\"failed to get owner by domain-" + domain + "\"}"
		return nil, errors.New(jsonResp)
	}
	return owner, nil
}

// GetDomainsByOwner 通过用户名获取域名集合
func GetDomainsByOwner(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	owner := args[0]
	domainList, err := stub.GetState(owner)
	if err != nil {
		jsonResp := "{\"Error\":\"failed to get domain list by owner" + owner + "\"}"
		return nil, errors.New(jsonResp)
	}
	return domainList, nil
}

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
