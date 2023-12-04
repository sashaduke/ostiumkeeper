package main

import (
	"encoding/json"
	"github.com/ethereum/go-ethereum/common"
	"net/http"
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

	// TODO: UPDATE THIS
	contractAddress := common.HexToAddress("0xMyContractAddress")
	_, err = NewSepoliaContract(client, contractAddress)
	if err != nil {
		http.Error(w, "Failed to instantiate contract", http.StatusInternalServerError)
		return
	}
}
