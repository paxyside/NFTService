package contract

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"log/slog"
	"math/big"
	"nft_service/infrastructure/config"
	"strings"
	"sync"
	"time"
)

type NFTContract struct {
	client       *ethclient.Client
	cfg          *config.Config
	contractABI  string
	cache        *big.Int
	cacheUpdated int64
	mu           sync.RWMutex
}

func NewNFTContract(url string, cfg *config.Config, contractABI string) (*NFTContract, error) {
	client, err := ethclient.Dial(url)
	if err != nil {
		return nil, err
	}

	contract := &NFTContract{
		client:      client,
		cfg:         cfg,
		contractABI: contractABI,
	}

	return contract, nil
}

func (m *NFTContract) Mint(owner string, uniqueHash string, mediaURL string) (string, error) {

	var (
		l         = slog.Default()
		startTime = time.Now()
	)

	parsedAbi, err := abi.JSON(strings.NewReader(m.contractABI))
	if err != nil {
		return "", fmt.Errorf("failed to parse contract ABI: %w", err)
	}

	txData, err := parsedAbi.Pack("mint", common.HexToAddress(owner), uniqueHash, mediaURL)
	if err != nil {
		return "", fmt.Errorf("failed to pack mint transaction data: %w", err)
	}

	nonce, err := m.client.PendingNonceAt(context.Background(), common.HexToAddress(m.cfg.UserAddress))
	if err != nil {
		return "", fmt.Errorf("failed to get pending nonce: %w", err)
	}

	gasPrice, err := m.client.SuggestGasPrice(context.Background())
	if err != nil {
		return "", fmt.Errorf("failed to suggest gas price: %w", err)
	}

	toAddress := common.HexToAddress(m.cfg.ContractAddress)
	value := big.NewInt(0)
	gasLimit := uint64(300000)

	unsignedTx := types.NewTx(&types.DynamicFeeTx{
		ChainID:   big.NewInt(m.cfg.ChainID),
		Nonce:     nonce,
		GasTipCap: gasPrice,
		GasFeeCap: gasPrice,
		Gas:       gasLimit,
		To:        &toAddress,
		Value:     value,
		Data:      txData,
	})

	privateKeyBytes, err := crypto.HexToECDSA(m.cfg.UserPrivateKey)
	if err != nil {
		return "", fmt.Errorf("failed to convert private key: %w", err)
	}

	signedTx, err := types.SignTx(unsignedTx, types.LatestSignerForChainID(big.NewInt(m.cfg.ChainID)), privateKeyBytes)
	if err != nil {
		return "", fmt.Errorf("failed to sign transaction: %w", err)
	}

	err = m.client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return "", fmt.Errorf("failed to send transaction: %w", err)
	}

	latency := time.Now().Sub(startTime).Milliseconds()
	l.Info("mint transaction sent", slog.Float64("latency", float64(latency)*0.001))

	return signedTx.Hash().Hex(), nil
}

// getTotalSupplyFromContract returns total supply from contract
func (m *NFTContract) getTotalSupplyFromContract() (*big.Int, error) {
	var (
		toAddress   = common.HexToAddress(m.cfg.ContractAddress)
		totalSupply *big.Int
	)

	parsedAbi, err := abi.JSON(strings.NewReader(m.contractABI))
	if err != nil {
		return nil, fmt.Errorf("failed to parse contract ABI: %w", err)
	}

	callData, err := parsedAbi.Pack("totalSupply")
	if err != nil {
		return nil, fmt.Errorf("failed to pack totalSupply call data: %w", err)
	}

	msg := ethereum.CallMsg{
		To:   &toAddress,
		Data: callData,
	}

	result, err := m.client.CallContract(context.Background(), msg, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to call totalSupply: %w", err)
	}

	err = parsedAbi.UnpackIntoInterface(&totalSupply, "totalSupply", result)
	if err != nil {
		return nil, fmt.Errorf("failed to unpack totalSupply result: %w", err)
	}

	return totalSupply, nil
}

// TotalSupply returning total supply from cache (not exact, but more efficient)
func (m *NFTContract) TotalSupply() (*big.Int, error) {
	m.mu.RLock()
	cache := m.cache
	cacheUpdated := m.cacheUpdated
	m.mu.RUnlock()

	if cache != nil && time.Now().Unix()-cacheUpdated < 60 {
		return cache, nil
	}
	totalSupply, err := m.updateTotalSupplyCache()
	if err != nil {
		return nil, err
	}

	return totalSupply, nil
}

// ExactTotalSupply returning exact total supply
func (m *NFTContract) ExactTotalSupply() (*big.Int, error) {
	return m.getTotalSupplyFromContract()
}

// updateTotalSupplyCache update cache with total supply
func (m *NFTContract) updateTotalSupplyCache() (*big.Int, error) {
	totalSupply, err := m.getTotalSupplyFromContract()
	if err != nil {
		return nil, err
	}
	m.mu.Lock()
	m.cache = totalSupply
	m.cacheUpdated = time.Now().Unix()
	m.mu.Unlock()

	return totalSupply, nil
}
