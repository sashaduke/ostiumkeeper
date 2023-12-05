package main

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stretchr/testify/require"
)

func TestConnectToEthereum(t *testing.T) {
	client, auth, err := connectToEthereum()
	require.Nil(t, err)
	require.NotNil(t, client)
	require.NotNil(t, auth)
}

func TestInstantiateContract(t *testing.T) {
	client, _, _ := connectToEthereum()
	contractAddress := common.HexToAddress(storageContractAddress)
	contract, err := instantiateContract(client, contractAddress)
	require.Nil(t, err)
	require.NotNil(t, contract)
}

func connectToLocalEthereum() (*ethclient.Client, error) {
	client, err := ethclient.Dial("http://localhost:8545") // Ganache default RPC port
	if err != nil {
		return nil, err
	}
	return client, nil
}

func TestLocalEthereumConnection(t *testing.T) {
	client, err := connectToLocalEthereum()
	require.Nil(t, err)
	require.NotNil(t, client)
}
