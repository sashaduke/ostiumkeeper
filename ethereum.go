package main

import (
	"context"
	"math/big"
	"os"
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
	client, auth, err := connectToEthereum()
	if err != nil {
		logger.Fatalf("ethereum connection error: %v\n", err)
	}

	contractAddress := common.HexToAddress(storageContractAddress)
	contract, err := instantiateContract(client, contractAddress)
	if err != nil {
		logger.Fatalf("contract instantiation error: %v\n", err)
	}

	for {
		// Write to blockchain every 15s to minimise gas costs and RPC requests
		time.Sleep(15 * time.Second)
		if err := writeToContract(client, auth, contract); err != nil {
			logger.Fatalf("contract interaction error: %v\n", err)
		}
	}
}

// Updates the smart contract with the latest timestamped data.
func writeToContract(client *ethclient.Client, auth *bind.TransactOpts, contract *Storage) error {
	nonce, err := client.PendingNonceAt(context.Background(), auth.From)
	if err != nil {
		return err
	}
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return err
	}

	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)      // eth sent in wei
	auth.GasLimit = uint64(3000000) // in gas units
	auth.GasPrice = gasPrice

	data, err := retrieveDataRedis()
	if err != nil {
		return err
	}

	tx, err := contract.Store(auth, data.Value)
	if err != nil {
		return err
	}
	logger.Printf("\nPrice update sent to blockchain! Tx hash:\n%s\n\n", tx.Hash().Hex())
	return nil
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
