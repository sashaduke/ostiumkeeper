package main

import (
	"encoding/json"
	"net/http"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
)

// handleData handles the "/data" HTTP route.
func handleData(w http.ResponseWriter, r *http.Request) {
	data, err := retrieveDataRedis()
	if err != nil {
		respondWithError(w, "failed to retrieve data from redis")
		return
	}

	if err := json.NewEncoder(w).Encode(data); err != nil {
		respondWithError(w, "failed to encode data")
	}
}

// handleContracts handles the "/contracts" HTTP route.
func handleContracts(w http.ResponseWriter, r *http.Request) {
	client, _, err := connectToEthereum()
	if err != nil {
		respondWithError(w, "failed to connect to Ethereum client")
		return
	}

	contractAddress := common.HexToAddress(storageContractAddress)
	contract, err := instantiateContract(client, contractAddress)
	if err != nil {
		respondWithError(w, "failed to instantiate contract")
		return
	}

	data, err := contract.Retrieve(&bind.CallOpts{})
	if err != nil {
		respondWithError(w, "failed to retrieve contract storage data")
		return
	}

	if err := json.NewEncoder(w).Encode(data); err != nil {
		respondWithError(w, "failed to encode data")
	}
}

func respondWithError(w http.ResponseWriter, message string) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(message))
	logger.Fatalf("handler error: %s", message)
}
