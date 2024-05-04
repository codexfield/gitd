package storage

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

var (
	EnvChainID    = "GREENFIELD_CHAIN_ID"
	EnvPrivateKey = "GREENFIELD_PRIVATE_KEY"
)

func InitTestStorage() (*GnfdStorage, error) {
	chainID := os.Getenv(EnvChainID)
	if chainID == "" {
		panic(fmt.Sprintf("Please set the environment variable: %s", EnvChainID))
	}

	privateKey := os.Getenv(EnvPrivateKey)
	if privateKey == "" {
		panic(fmt.Sprintf("Please set the enviroment variable: %s", EnvPrivateKey))
	}
	rpcAddress := "https://greenfield-chain-ap.bnbchain.org:443"
	bucketName := "codex-4-gitd-test"
	return NewStorage(chainID, rpcAddress, privateKey, bucketName)
}

func TestStorage_Put(t *testing.T) {
	store, err := InitTestStorage()
	assert.NoError(t, err)
	//err = store.put("test2", []byte("test-value"), false)
	//assert.NoError(t, err)
	_, err = store.get("refs/refs/heads/main\n")
	assert.NoError(t, err)
}
