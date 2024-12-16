package contract

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"math/big"
	"strings"
	"time"
)

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

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := m.client.CallContract(ctx, msg, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to call totalSupply: %w", err)
	}

	err = parsedAbi.UnpackIntoInterface(&totalSupply, "totalSupply", result)
	if err != nil {
		return nil, fmt.Errorf("failed to unpack totalSupply result: %w", err)
	}

	return totalSupply, nil
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
