package config

import (
	"errors"
	"log/slog"
	"os"
	"strconv"
	"time"
)

type Config struct {
	Host                string
	Port                string
	DBURI               string
	CacheUpdateInterval time.Duration
	UserAddress         string
	UserPrivateKey      string
	NetworkName         string
	InfuraApiKey        string
	ChainID             int64
	ContractAddress     string
	ContractABIPath     string
}

func LoadConfig() (*Config, error) {

	var l = slog.Default()

	host := os.Getenv("HOST")
	if host == "" {
		l.Error("HOST is not set")
		return nil, errors.New("HOST is not set")
	}

	port := os.Getenv("PORT")
	if port == "" {
		l.Error("PORT is not set")
		return nil, errors.New("PORT is not set")
	}

	dbURI := os.Getenv("DB_URI")
	if dbURI == "" {
		l.Error("DB_URI is not set")
		return nil, errors.New("DB_URI is not set")
	}

	cacheUpdateInterval := os.Getenv("CACHE_UPDATE_INTERVAL")
	if cacheUpdateInterval == "" {
		l.Error("CACHE_UPDATE_DURATION is not set")
		return nil, errors.New("CACHE_UPDATE_DURATION is not set")
	}

	intCacheUpdateInterval, err := strconv.ParseInt(cacheUpdateInterval, 10, 64)
	if err != nil {
		l.Error("CACHE_UPDATE_DURATION is not integer", "error", err)
		return nil, errors.New("CACHE_UPDATE_DURATION is not integer")
	}

	userAddress := os.Getenv("USER_ADDRESS")
	if userAddress == "" {
		l.Error("USER_ADDRESS is not set")
		return nil, errors.New("USER_ADDRESS is not set")
	}

	userPrivateKey := os.Getenv("USER_PRIVATE_KEY")
	if userPrivateKey == "" {
		l.Error("USER_PRIVATE_KEY is not set")
		return nil, errors.New("USER_PRIVATE_KEY is not set")
	}

	networkName := os.Getenv("NETWORK_NAME")
	if networkName == "" {
		l.Error("NETWORK_NAME is not set")
		return nil, errors.New("NETWORK_NAME is not set")
	}

	infuraApiKey := os.Getenv("INFURA_API_KEY")
	if infuraApiKey == "" {
		l.Error("INFURA_API_KEY is not set")
		return nil, errors.New("INFURA_API_KEY is not set")
	}

	chainID := os.Getenv("CHAIN_ID")
	if chainID == "" {
		l.Error("CHAIN_ID is not set")
		return nil, errors.New("CHAIN_ID is not set")
	}

	intChainID, err := strconv.ParseInt(chainID, 10, 64)
	if err != nil {
		l.Error("failed to parse chain ID", "error", err)
		return nil, err
	}

	contractAddress := os.Getenv("CONTRACT_ADDRESS")
	if contractAddress == "" {
		l.Error("CONTRACT_ADDRESS is not set")
		return nil, errors.New("CONTRACT_ADDRESS is not set")
	}
	contractABIPath := os.Getenv("CONTRACT_ABI_PATH")
	if contractABIPath == "" {
		l.Error("CONTRACT_ABI_PATH is not set")
		return nil, errors.New("CONTRACT_ABI_PATH is not set")
	}

	return &Config{
		Host:                host,
		Port:                port,
		DBURI:               dbURI,
		CacheUpdateInterval: time.Duration(intCacheUpdateInterval) * time.Second,
		UserAddress:         userAddress,
		UserPrivateKey:      userPrivateKey,
		NetworkName:         networkName,
		InfuraApiKey:        infuraApiKey,
		ChainID:             intChainID,
		ContractAddress:     contractAddress,
		ContractABIPath:     contractABIPath,
	}, nil
}
