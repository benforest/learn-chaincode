/*
Copyright IBM Corp 2016 All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

		 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"errors"
	"fmt"
	"encoding/json"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

type Account struct{
	Name string  `json:"name"`
	Role string `json:"role"`
	Balance 		float64 `json:"balance"`
}

type Transaction struct {
	FromName string   `json:"fromName"`
	ToName   string   `json:"toName"`
	Quantity float64 `json:"quantity"`
}
// ============================================================================================================================
// Main
// ============================================================================================================================
func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

// Init resets all the things
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}

	return nil, nil
}

// Invoke is our entry point to invoke a chaincode function
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("invoke is running " + function)

	// Handle different functions
	if function == "init" {													//initialize the chaincode state, used as reset
		return t.Init(stub, "init", args)
	}
	fmt.Println("Invoke running. Function: " + function)

	if function == "transfer" {
		return t.transfer(stub, args)
	}  else if function == "createAccount" {
		return t.createAccount(stub, args)
	}

	fmt.Println("invoke did not find func: " + function)					//error

	return nil, errors.New("Received unknown function invocation: " + function)
}

func (t *SimpleChaincode) createAccount(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	/*		0
		json
	  	{
			  "name":"",
			  "role":"",
			  "balance": ""
		}
	*/
	//need one argfmt.Println("Creating account")
	if len(args) != 1 {
		fmt.Println("Error obtaining user")
		return nil, errors.New("createAccount accepts a Account argument")
	}

	var account Account

	fmt.Println("Unmarshalling Account")
	err := json.Unmarshal([]byte(args[0]), &account)
	if err != nil {
			fmt.Println("Error Unmarshalling Account")
			return nil, errors.New("Invalid Account")
	}

	accountBytes, err := json.Marshal(&account)
	fmt.Println("initializing account.")
		err = stub.PutState(account.Name, accountBytes)

		if err == nil {
			fmt.Println("created account" + account.Name)
			return nil, nil
		} else {
			fmt.Println("failed to create initialize account for " + account.Name)
			return nil, errors.New("failed to initialize an account for " + account.Name + " => " + err.Error())
		}
}

func (t *SimpleChaincode) transfer(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	fmt.Println("Transferring Money")
	/*		0
		json
	  	{
			  "fromName":"",
			  "toName":"",
			  "quantity": ""
		}
	*/
	//need one arg
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting transaction")
	}

	var tr Transaction

	fmt.Println("Unmarshalling Transaction")
	err := json.Unmarshal([]byte(args[0]), &tr)
	if err != nil {
			fmt.Println("Error Unmarshalling Transaction")
			return nil, errors.New("Invalid Transaction")
	}

	var fromAccount Account
	fmt.Println("Getting State on sender " + tr.FromName)
	fromBytes, err := stub.GetState(tr.FromName)
	if err != nil {
		fmt.Println("fromName not found")
		return nil, errors.New("fromName not found" + tr.FromName)
	}
	fmt.Println("Unmarshalling FromAccount ")
	err = json.Unmarshal(fromBytes, &fromAccount)
	if err != nil {
		fmt.Println("Error unmarshalling account " + tr.FromName)
		return nil, errors.New("Error unmarshalling account " + tr.FromName)
	}

	var toAccount Account
	fmt.Println("Getting State on ToAccount " + tr.ToName)
	toBytes, err := stub.GetState(tr.ToName)
	if err != nil {
		fmt.Println("Account not found " + tr.ToName)
		return nil, errors.New("Account not found " + tr.ToName)
	}
	fmt.Println("Unmarshalling ToAccount")
	err = json.Unmarshal(toBytes, &toAccount)
	if err != nil {
		fmt.Println("Error unmarshalling account " + tr.ToName)
		return nil, errors.New("Error unmarshalling account " + tr.ToName)
	}


	if fromAccount.Balance < tr.Quantity {
		fmt.Println("The FromAccount " + tr.FromName + "doesn't have enough to transfer")
		return nil, errors.New("The FromAccount " + tr.FromName + "doesn't have enough to transfer")
	} else {
		fmt.Println("The FromAccount has enough to be transferred")
	}

	toAccount.Balance += tr.Quantity
	fromAccount.Balance -= tr.Quantity

	fromBytesToWrite, err := json.Marshal(&fromAccount)
		if err != nil {
			fmt.Println("Error marshalling the fromAccount")
			return nil, errors.New("Error marshalling the fromAccount")
		}
		fmt.Println("Put state on fromAccount")
		err = stub.PutState(tr.FromName, fromBytesToWrite)
		if err != nil {
			fmt.Println("Error writing the fromAccount back")
			return nil, errors.New("Error writing the fromAccount back")
	}

	toBytesToWrite, err := json.Marshal(&toAccount)
	if err != nil {
		fmt.Println("Error marshalling the toAccount")
		return nil, errors.New("Error marshalling the toAccount")
	}
	fmt.Println("Put state on toAccount")
	err = stub.PutState(tr.ToName, toBytesToWrite)
	if err != nil {
		fmt.Println("Error writing the toAccount back")
		return nil, errors.New("Error writing the toAccount back")
	}

	transactionToWrite, err := json.Marshal(&tr)
	if err != nil {
		fmt.Println("Error marshalling the transaction")
		return nil, errors.New("Error marshalling the transaction")
	}
	fmt.Println("Put state on transaction")
	err = stub.PutState(tr.FromName, transactionToWrite)
	if err != nil {
		fmt.Println("Error writing the transaction back")
		return nil, errors.New("Error writing the transaction back")
	}

	fmt.Println("Successfully completed Invoke")
	return nil, nil
}

// Query is our entry point for queries
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("query is running " + function)

	// Handle different functions
	if function == "getAccount" {
	 fmt.Println("Getting particular account")
	 account, err := t.getAccount(args[0], stub)
	 if err != nil {
		 fmt.Println("Error Getting particular account")
		 return nil, err
	 } else {
		 accountBytes, err1 := json.Marshal(&account)
		 if err1 != nil {
			 fmt.Println("Error marshalling the account")
			 return nil, err1
		 }
		 fmt.Println("All success, returning the account")
		 return accountBytes, nil
	 }
 }

	 if function == "getTransaction" {
 	 fmt.Println("Getting particular transaction")
 	 transaction, err := t.getTransaction(args[0], stub)
 	 if err != nil {
 		 fmt.Println("Error Getting particular transaction")
 		 return nil, err
 	 } else {
 		 transactionBytes, err1 := json.Marshal(&transaction)
 		 if err1 != nil {
 			 fmt.Println("Error marshalling the transaction")
 			 return nil, err1
 		 }
 		 fmt.Println("All success, returning the transaction")
 		 return transactionBytes, nil
 	 }
 }
	return nil, errors.New("Received unknown function query: " + function)
}

func (t *SimpleChaincode) getAccount(accountName string, stub shim.ChaincodeStubInterface) (Account, error) {
	var account Account

	accountBytes, err := stub.GetState(accountName)
	if err != nil {
		fmt.Println("Error retrieving account " + accountName)
		return account, errors.New("Error retrieving account " + accountName)
	}

	err = json.Unmarshal(accountBytes, &account)
	if err != nil {
		fmt.Println("Error unmarshalling account " + accountName)
		return account, err
	}

	return account, nil
}

func (t *SimpleChaincode) getTransaction(transactionName string, stub shim.ChaincodeStubInterface) (Transaction, error) {
	var transaction Transaction

	transactionBytes, err := stub.GetState(transactionName)
	if err != nil {
		fmt.Println("Error retrieving account " + transactionName)
		return transaction, errors.New("Error retrieving account " + transactionName)
	}

	err = json.Unmarshal(transactionBytes, &transaction)
	if err != nil {
		fmt.Println("Error unmarshalling account " + transactionName)
		return transaction, errors.New("Error unmarshalling account " + transactionName)
	}

	return transaction, nil
}
