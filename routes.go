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
		respondWithError(w, "failed to retrieve data from redis", err, http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(data); err != nil {
		respondWithError(w, "failed to JSON encode data", err, http.StatusInternalServerError)
	}
}

// handleContracts handles the "/contracts" HTTP route.
func handleContracts(w http.ResponseWriter, r *http.Request) {
	// Will try connect 5 times, 3s apart, before failing
	client, _, err := connectToEthereumWithRetry(5)
	if err != nil {
		respondWithError(w, "failed to connect to blockchain", err, http.StatusInternalServerError)
		return
	}

	contractAddress := common.HexToAddress(storageContractAddress)
	contract, err := instantiateContract(client, contractAddress)
	if err != nil {
		respondWithError(w, "failed to instantiate contract object", err, http.StatusInternalServerError)
		return
	}

	data, err := contract.Retrieve(&bind.CallOpts{})
	if err != nil {
		respondWithError(w, "failed to retrieve contract storage data", err, http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(data); err != nil {
		respondWithError(w, "failed to JSON encode data", err, http.StatusInternalServerError)
	}
}

func respondWithError(w http.ResponseWriter, message string, err error, statusCode int) {
	w.WriteHeader(statusCode)
	if _, e := w.Write([]byte(message)); e != nil {
		logger.Printf("couldn't write error response: %v", e)
	}
	logger.Printf("API route handler error, %s:\n%v", message, err)
}
