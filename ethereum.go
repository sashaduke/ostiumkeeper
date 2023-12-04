package main

import (
	"context"
	"log"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

const (
	SepoliaRPCEndpoint = "https://sepolia.infura.io/v3/131bd995e0764b2da6be91ee9058dc91"
	ECDSAPrivkeyHex    = "fe05041e74295604ff8f76dc24847c06e93c015da608b4281446c7de6f54cc46"
)

// Keeper logic to interact with Ethereum blockchain.
func keeper() {
	client, auth, err := connectToEthereum()
	if err != nil {
		log.Fatalf("ethereum connection error: %v\n", err)
	}

	for {
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
		auth.GasLimit = uint64(3000000) // in units
		auth.GasPrice = gasPrice

		data, err := retrieveDataRedis()
		if err != nil {
			log.Printf("redis get error: %v\n", err)
			continue
		}
		log.Println(data)

		//TODO: UPDATE THIS
		contractAddress := common.HexToAddress("0xMyContractAddress")
		_, err = NewSepoliaContract(client, contractAddress)
		if err != nil {
			log.Printf("contract instantiation error: %v\n", err)
			continue
		}

		time.Sleep(time.Second)
	}
}

// connectToEthereum establishes a connection to an Ethereum client and creates an authenticated session.
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

// NewSepoliaContract instantiates a new smart contract.
func NewSepoliaContract(client *ethclient.Client, address common.Address) (*bind.BoundContract, error) {
	contractABIJSON := `{"constant":false,"inputs":[],"name":"getLatestData","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"nonpayable","type":"function"}` //TODO: Replace with actual contract's ABI JSON
	parsedABI, err := abi.JSON(strings.NewReader(contractABIJSON))
	if err != nil {
		return nil, err
	}

	return bind.NewBoundContract(address, parsedABI, client, client, client), nil
}
