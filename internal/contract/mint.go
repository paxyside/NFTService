package contract

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"log/slog"
	"math/big"
	"nft_service/internal/domain"
	"time"
)

func (m *NFTContract) Mint(token *domain.Token) (*domain.Token, error) {

	var (
		l         = slog.Default()
		startTime = time.Now()
	)

	txData, err := m.parsedABI.Pack("mint", common.HexToAddress(token.Owner), token.UniqueHash, token.MediaUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to pack mint transaction data: %w", err)
	}

	nonce, err := m.client.PendingNonceAt(context.Background(), common.HexToAddress(m.cfg.UserAddress))
	if err != nil {
		return nil, fmt.Errorf("failed to get pending nonce: %w", err)
	}

	gasPrice, err := m.client.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to suggest gas price: %w", err)
	}

	toAddress := common.HexToAddress(m.cfg.ContractAddress)

	unsignedTx := types.NewTx(&types.DynamicFeeTx{
		ChainID:   big.NewInt(m.cfg.ChainID),
		Nonce:     nonce,
		GasTipCap: gasPrice,
		GasFeeCap: gasPrice,
		Gas:       uint64(300000),
		To:        &toAddress,
		Value:     big.NewInt(0),
		Data:      txData,
	})

	privateKeyBytes, err := crypto.HexToECDSA(m.cfg.UserPrivateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to convert private key: %w", err)
	}

	signedTx, err := types.SignTx(unsignedTx, types.LatestSignerForChainID(big.NewInt(m.cfg.ChainID)), privateKeyBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to sign transaction: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	err = m.client.SendTransaction(ctx, signedTx)
	if err != nil {
		return nil, fmt.Errorf("failed to send transaction: %w", err)
	}

	token.TxHash = signedTx.Hash().Hex()

	latency1 := time.Now().Sub(startTime).Milliseconds()
	l.Info("mint transaction sent to contract", slog.Float64("latency", float64(latency1)*0.001))

	return token, nil
}
