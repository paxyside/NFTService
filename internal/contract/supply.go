package contract

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"log/slog"
	"math/big"
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

	callData, err := m.parsedABI.Pack("totalSupply")
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

	err = m.parsedABI.UnpackIntoInterface(&totalSupply, "totalSupply", result)
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

func (m *NFTContract) StartCacheUpdater(ctx context.Context, interval time.Duration) {
	l := slog.Default()
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			_, err := m.updateTotalSupplyCache()
			if err != nil {
				l.Error("failed to update total supply cache", slog.Any("error", err))
			}
		case <-ctx.Done():
			l.Info("cache updater stopped")
			return
		}
	}
}
