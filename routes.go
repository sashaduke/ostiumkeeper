package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
)

// handleData handles the "/data" HTTP route.
func handleData(w http.ResponseWriter, r *http.Request) {
	data, err := retrieveDataRedis()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err = json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "Error encoding data", http.StatusInternalServerError)
	}
}

// handleContracts handles the "/contracts" HTTP route.
func handleContracts(w http.ResponseWriter, r *http.Request) {
	client, _, err := connectToEthereum()
	if err != nil {
		http.Error(w, "Failed to connect to Ethereum client", http.StatusInternalServerError)
		return
	}

	contractAddress := common.HexToAddress("48eB2302cfEc7049820b66FC91955C5d250b3fF9")
	contract, err := instantiateContract(client, contractAddress)
	if err != nil {
		http.Error(w, "Failed to instantiate contract", http.StatusInternalServerError)
		return
	}

	data, err := contract.Retrieve(&bind.CallOpts{})
	if err != nil {
		log.Fatalf("contract call error: %v\n", err)
	}

	if err = json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "Error encoding data", http.StatusInternalServerError)
	}
}
