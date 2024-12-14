package utils

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateUniqueHash(t *testing.T) {
	hash, err := GenerateUniqueHash()
	assert.NoError(t, err)
	assert.NotEmpty(t, hash)
	assert.Equal(t, 20, len(hash))
}

func TestGenerateInfuraURL(t *testing.T) {
	networkName := "mainnet"
	apiKey := "testapikey"

	url, err := GenerateInfuraURL(networkName, apiKey)
	assert.NoError(t, err)
	assert.Equal(t, "https://mainnet.infura.io/v3/testapikey", url)

	_, err = GenerateInfuraURL("", apiKey)
	assert.NoError(t, err)
}

func TestLoadABIFromFile(t *testing.T) {
	filePath := "test_abi.json"
	abiContent := `{"abi": "test"}`
	err := os.WriteFile(filePath, []byte(abiContent), 0644)
	assert.NoError(t, err)

	abi, err := LoadABIFromFile(filePath)
	assert.NoError(t, err)
	assert.Equal(t, abiContent, abi)

	err = os.Remove(filePath)
	assert.NoError(t, err)

	_, err = LoadABIFromFile(filePath)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to read ABI file")
}
