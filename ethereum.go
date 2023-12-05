package main

import (
	"context"
	"math/big"
	"os"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

var (
	storageContractAddress = getEnv("CONTRACT_ADDRESS", "48eB2302cfEc7049820b66FC91955C5d250b3fF9")
	blockchainRPCEndpoint  = getEnv("RPC_ENDPOINT", "https://sepolia.infura.io/v3/131bd995e0764b2da6be91ee9058dc91")
	privkeyHexECDSA        = getEnv("PRIVKEY_HEX", "fe05041e74295604ff8f76dc24847c06e93c015da608b4281446c7de6f54cc46")
	contractWriteFrequency = getEnv("CONTRACT_WRITE_FREQ", "15")
)

// Get environment variables or fallback to above values.
func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

// Keeper logic to interact with EVM-based blockchain.
func keeper() {
	client, auth, err := connectToEthereumWithRetry(5)
	if err != nil {
		logger.Fatalf("ethereum connection error: %v\n", err)
	}

	contractAddress := common.HexToAddress(storageContractAddress)
	contract, err := instantiateContract(client, contractAddress)
	if err != nil {
		logger.Fatalf("contract instantiation error: %v\n", err)
	}

	latestUpdate, err := retrieveDataRedis()
	if err != nil {
		latestUpdate = Data{Timestamp: time.Now().UTC()}
	}

	// Write to blockchain no more than every 15s by default to minimise gas costs and RPC requests
	var freq time.Duration
	if freqInt, err := strconv.Atoi(contractWriteFrequency); err != nil {
		freq = time.Second * 15
	} else {
		freq = time.Second * time.Duration(freqInt)
	}

	for {
		time.Sleep(3 * time.Second)
		data, err := retrieveDataRedis()
		// Check that the most recent timestamp is indeed 15s+ later than the last one
		if data.Timestamp.Before(latestUpdate.Timestamp.Add(freq)) ||
			data.Timestamp.Equal(latestUpdate.Timestamp.Add(freq)) || err != nil {
			continue
		}
		if err := writeToContract(client, auth, contract, data); err != nil {
			logger.Fatalf("contract interaction error: %v\n", err)
		}
		latestUpdate = data
	}
}

// Updates the smart contract with the latest timestamped data.
func writeToContract(client *ethclient.Client, auth *bind.TransactOpts, contract *Storage, data Data) error {
	nonce, err := client.PendingNonceAt(context.Background(), auth.From)
	if err != nil {
		return err
	}
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return err
	}

	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)      // eth to send in wei
	auth.GasLimit = uint64(3000000) // gas limit in gas units
	auth.GasPrice = gasPrice

	tx, err := contract.Store(auth, data.Value)
	if err != nil {
		return err
	}
	logger.Printf("\nPrice update contract call tx broadcasted to blockchain! Tx hash:\n%s\n\n", tx.Hash().Hex())
	return nil
}

func connectToEthereumWithRetry(maxRetries int) (*ethclient.Client, *bind.TransactOpts, error) {
	var client *ethclient.Client
	var auth *bind.TransactOpts
	var err error

	for attempt := 0; attempt < maxRetries; attempt++ {
		client, auth, err = connectToEthereum()
		if err == nil {
			return client, auth, nil
		}
		time.Sleep(3 * time.Second) // Wait for 3 seconds before retrying
	}
	return nil, nil, err
}

// connectToEthereum establishes a connection to an Ethereum client.
func connectToEthereum() (*ethclient.Client, *bind.TransactOpts, error) {
	client, err := ethclient.Dial(blockchainRPCEndpoint)
	if err != nil {
		return nil, nil, err
	}

	privateKey, err := crypto.HexToECDSA(privkeyHexECDSA)
	if err != nil {
		return nil, nil, err
	}

	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		return nil, nil, err
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		return nil, nil, err
	}

	return client, auth, nil
}

// instantiateContract creates an object with contract methods using generated ABI bindings.
func instantiateContract(client *ethclient.Client, address common.Address) (*Storage, error) {
	contract, err := NewStorage(address, client)
	if err != nil {
		return nil, err
	}
	return contract, nil
}
