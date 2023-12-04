package main

import (
	"context"
	"log"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

const (
	StorageContractAddress = "48eB2302cfEc7049820b66FC91955C5d250b3fF9"
	SepoliaRPCEndpoint     = "https://sepolia.infura.io/v3/131bd995e0764b2da6be91ee9058dc91"
	ECDSAPrivkeyHex        = "fe05041e74295604ff8f76dc24847c06e93c015da608b4281446c7de6f54cc46"
)

// Keeper logic to interact with Ethereum blockchain (Sepolia Testnet)
func keeper() {
	client, auth, err := connectToEthereum()
	if err != nil {
		log.Fatalf("ethereum connection error: %v\n", err)
	}

	contractAddress := common.HexToAddress(StorageContractAddress)
	contract, err := instantiateContract(client, contractAddress)
	if err != nil {
		log.Fatalf("contract instantiation error: %v\n", err)
	}

	for {
		// Prevents us spending too much gas or getting rate-limited by RPC provider
		time.Sleep(20 * time.Second)

		nonce, err := client.PendingNonceAt(context.Background(), auth.From)
		if err != nil {
			log.Fatalf("nonce error: %v\n", err)
		}
		gasPrice, err := client.SuggestGasPrice(context.Background())
		if err != nil {
			log.Fatalf("gas price error: %v\n", err)
		}

		auth.Nonce = big.NewInt(int64(nonce))
		auth.Value = big.NewInt(0)      // in wei
		auth.GasLimit = uint64(3000000) // in gas units
		auth.GasPrice = gasPrice

		data, err := retrieveDataRedis()
		if err != nil {
			log.Fatalf("redis get error: %v\n", err)
		}

		tx, err := contract.Store(auth, data.Value)
		if err != nil {
			log.Fatalf("contract call error: %v\n", err)
		}
		log.Printf("\nPrice update sent to blockchain! Tx hash:\n%s\n\n", tx.Hash().Hex())
	}
}

// connectToEthereum establishes a connection to an Ethereum client and creates an authenticated session
func connectToEthereum() (*ethclient.Client, *bind.TransactOpts, error) {
	client, err := ethclient.Dial(SepoliaRPCEndpoint)
	if err != nil {
		return nil, nil, err
	}

	privateKey, err := crypto.HexToECDSA(ECDSAPrivkeyHex)
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

// instantiateContract creates an object with contract methods using generated ABI bindings
func instantiateContract(client *ethclient.Client, address common.Address) (*Storage, error) {
	contract, err := NewStorage(address, client)
	if err != nil {
		return nil, err
	}
	return contract, nil
}
