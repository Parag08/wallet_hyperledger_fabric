/*
 * Copyright Wavenet All Rights Reserved
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package main

import (
	"fmt"
	"crypto/sha256"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
	"encoding/json"
	"strings"
	"strconv"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

// Wallet implements a structure to hold wallet information
type wallet struct {
	Name string `json:"name"`
	Owner string `json:"owner"`
	Balance float64 `json:"balance"`
	Password string `json:"password"`
}

// Init is called during chaincode instantiation to initialize any
// data. Note that chaincode upgrade also calls this function to reset
// or to migrate data.
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface) peer.Response {
	// Get the args from the transaction proposal

	return shim.Success(nil)
}

// Invoke is called per transaction on the chaincode. Each transaction is
// either a 'get' or a 'set' on the asset created by Init function. The Set
// method may create a new asset by specifying a new key-value pair.
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	// Extract the function and args from the transaction proposal
	function, args := stub.GetFunctionAndParameters()
	fmt.Println("invoke is running " + function)

	if function == "initWallet" { //in it wallet
		return t.initWallet(stub, args)
	} else if function == "createWallet" {
		return t.createWallet(stub, args)
	} else if function == "transaction" {
		return t.transaction(stub, args)
	} else if function == "getWalletInfo" {
		return t.getWalletInfo(stub, args)
	}
	fmt.Println("invoke did not find func: " + function) //error
	return shim.Error("Received unknown function invocation")
}

func (t *SimpleChaincode) initWallet(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	masterPassword := "E7C038F31F4A241001E57D95DEF669F1E5DD2E40A46C22B15FE65CAC8C3DD03D"
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments expecting 1")
	}
	// finding hash of masterPassword provided
	sha_256 := sha256.New()
	sha_256.Write([]byte(args[0]))
	// check if wallet exist
	masterWalletAsBytes, err := stub.GetState("masterWallet")
	if err != nil {
		return shim.Error("Failed to get masterWallet: " + err.Error())
	} else if masterWalletAsBytes != nil {
		fmt.Println("masterWallet already exists")
		return shim.Error("masteWallet already exists")
	}
	sha_256_input := fmt.Sprintf("%X", sha_256.Sum(nil))
	fmt.Println(masterPassword,sha_256_input)
	if (strings.Compare(sha_256_input, masterPassword) == 0) {
		fmt.Println("master password verified \n intialising masterwallet")
		masterWallet := &wallet{"masterWallet","admin",100000,masterPassword}
		masterWalletJSONasBytes, err := json.Marshal(masterWallet)
		if err != nil {
			return shim.Error(err.Error())
		}
		// === Save wallet to state ===
		err = stub.PutState("masterWallet", masterWalletJSONasBytes)
		if err != nil {
			return shim.Error(err.Error())
		}
	} else  {
		return shim.Error("Wrong master password")
	}
	fmt.Println("- end init wallet")
	return shim.Success(nil)
}

func (t *SimpleChaincode) createWallet(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if (len(args) != 3) {
		return shim.Error("invalid length of arguments require 3 arguments [\"name of wallet\", \"name of owner\", \"password\"]")
	}
	if len(args[0]) <= 0 {
		return shim.Error("1st argument must be a non-empty string")
	}
	if len(args[1]) <= 0 {
		return shim.Error("2nd argument must be a non-empty string")
	}
	if len(args[2]) <= 0 {
		return shim.Error("3rd argument must be a non-empty string")
	}
	walletAsBytes, err := stub.GetState(args[0])
	if err != nil {
		return shim.Error("Failed to get Wallet: " + err.Error())
	} else if walletAsBytes != nil {
		fmt.Println("Wallet already exists")
		return shim.Error("Wallet already exists")
	}
	walletName := args[0]
	ownerName := args[1]
	// saving password in sha256 format
	sha_256 := sha256.New()
	sha_256.Write([]byte(args[2]))
	password := fmt.Sprintf("%X", sha_256.Sum(nil))
	walletData := &wallet{walletName, ownerName, 0, password}
	walletJSONasBytes, err := json.Marshal(walletData)
	if err != nil {
		return shim.Error(err.Error())
	}
	// === Save wallet to state ===
	err = stub.PutState(walletName, walletJSONasBytes)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(nil)
}

func (t *SimpleChaincode) transaction(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	// len of args should be 4
	// 0              1          2        4
	// fromWallet  towallet  amount   password
	if (len(args) != 4) {
		return shim.Error("invalid length of arguments require 3 arguments [\"name of wallet\", \"name of owner\", \"password\"]")
	}
	// checking sanity of input
	if len(args[0]) <= 0 {
		return shim.Error("1st argument must be a non-empty string")
	}
	if len(args[1]) <= 0 {
		return shim.Error("2nd argument must be a non-empty string")
	}
	if len(args[2]) <= 0 {
		return shim.Error("3rd argument must be a non-empty string")
	}
	if len(args[3]) <= 0 {
		return shim.Error("4th argument must be a non-empty string")
	}
	amount, err := strconv.ParseFloat(args[2],64)
	if err != nil {
		return shim.Error("3rd argument must be a numeric string")
	}
	if amount < 0 {
		return shim.Error("Amount should always be positive")
	}
	ownerWalletAsBytes, err := stub.GetState(args[0])
	if err != nil {
		return shim.Error("Failed to get Owner wallet:" + err.Error())
	} else if ownerWalletAsBytes == nil {
		return shim.Error("Owner wallet does not exist")
	}
	ownerWallet := wallet{}
	err = json.Unmarshal(ownerWalletAsBytes, &ownerWallet) //unmarshal it aka JSON.parse()
	if err != nil {
		return shim.Error(err.Error())
	}
	receiverWalletAsBytes, err := stub.GetState(args[1])
	if err != nil {
		return shim.Error("Failed to get Owner wallet:" + err.Error())
	} else if receiverWalletAsBytes == nil {
		return shim.Error("Owner wallet does not exist")
	}
	receiverWallet := wallet{}
	err = json.Unmarshal(receiverWalletAsBytes, &receiverWallet) //unmarshal it aka JSON.parse()
	if err != nil {
		return shim.Error(err.Error())
	}
	if t.Authenticate(ownerWallet, args[3]) {
		if ownerWallet.Balance >= amount {
			ownerWallet.Balance = ownerWallet.Balance - amount
			receiverWallet.Balance =  receiverWallet.Balance + amount

			// putting into current state of blockchain
			ownerWalletJSONasBytes, err := json.Marshal(ownerWallet)
			if err != nil {
				return shim.Error(err.Error())
			}
			// === Save wallet to state ===
			err = stub.PutState(ownerWallet.Name, ownerWalletJSONasBytes)
			if err != nil {
				return shim.Error(err.Error())
			}

			// putting into current state of blockchain
			receiverWalletJSONasBytes, err := json.Marshal(receiverWallet)
			if err != nil {
				return shim.Error(err.Error())
			}
			// === Save wallet to state ===
			err = stub.PutState(receiverWallet.Name, receiverWalletJSONasBytes)
			if err != nil {
				return shim.Error(err.Error())
			}
			return shim.Success(nil)
		} else {
			return shim.Error("You don't have sufficient balance for this transaction")
		}
	} else {
		return shim.Error("provide proper authentication")
	}
}

func convertStringToSHA256 (inputString string) string {
	sha_256 := sha256.New()
	sha_256.Write([]byte(inputString))
	return fmt.Sprintf("%X", sha_256.Sum(nil))
}

func (t *SimpleChaincode) Authenticate (walletInfo wallet,password string) bool {
	if (walletInfo.Password == convertStringToSHA256(password)) {
		return true
	}
	return false
}
func (t *SimpleChaincode) getWalletInfo(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}
	if len(args[0]) <= 0 {
		return shim.Error("1st argument must be a non-empty string")
	}
	if len(args[1]) <= 0 {
		return shim.Error("2nd argument must be a non-empty string")
	}
	ownerWalletAsBytes, err := stub.GetState(args[0])
	if err != nil {
		return shim.Error("Failed to get Owner wallet:" + err.Error())
	} else if ownerWalletAsBytes == nil {
		return shim.Error("Owner wallet does not exist")
	}
	ownerWallet := wallet{}
	err = json.Unmarshal(ownerWalletAsBytes, &ownerWallet) //unmarshal it aka JSON.parse()
	if err != nil {
		return shim.Error(err.Error())
	}
	if t.Authenticate(ownerWallet, args[1]) {
		return shim.Success(ownerWalletAsBytes)
	} else {
		return shim.Error("password is wrong")
	}
}

// main function starts up the chaincode in the container during instantiate
func main() {
	if err := shim.Start(new(SimpleChaincode)); err != nil {
		fmt.Printf("Error starting SimpleAsset chaincode: %s", err)
	}
}

