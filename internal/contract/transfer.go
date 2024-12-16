package contract

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"log/slog"
	"math/big"
	"strings"
	"time"
)

func (m *NFTContract) TransferToken(from string, to string, tokenId *big.Int) (string, error) {
	var (
		ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
		l           = slog.Default()
		startTime   = time.Now()
	)
	defer cancel()

	parsedAbi, err := abi.JSON(strings.NewReader(m.contractABI))
	if err != nil {
		return "", fmt.Errorf("failed to parse contract ABI: %w", err)
	}

	txData, err := parsedAbi.Pack("safeTransferFrom", common.HexToAddress(from), common.HexToAddress(to), tokenId)
	if err != nil {
		return "", fmt.Errorf("failed to pack transfer data: %w", err)
	}

	nonce, err := m.client.PendingNonceAt(ctx, common.HexToAddress(m.cfg.UserAddress))
	if err != nil {
		return "", fmt.Errorf("failed to get pending nonce: %w", err)
	}

	gasPrice, err := m.client.SuggestGasPrice(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to suggest gas price: %w", err)
	}

	toAddress := common.HexToAddress(m.cfg.ContractAddress)
	gasLimit := uint64(300000)

	unsignedTx := types.NewTx(&types.DynamicFeeTx{
		ChainID:   big.NewInt(m.cfg.ChainID),
		Nonce:     nonce,
		GasTipCap: gasPrice,
		GasFeeCap: gasPrice,
		Gas:       gasLimit,
		To:        &toAddress,
		Value:     big.NewInt(0),
		Data:      txData,
	})

	privateKey, err := crypto.HexToECDSA(m.cfg.UserPrivateKey)
	if err != nil {
		return "", fmt.Errorf("failed to convert private key: %w", err)
	}

	signedTx, err := types.SignTx(unsignedTx, types.LatestSignerForChainID(big.NewInt(m.cfg.ChainID)), privateKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign transaction: %w", err)
	}

	err = m.client.SendTransaction(ctx, signedTx)
	if err != nil {
		return "", fmt.Errorf("failed to send transaction: %w", err)
	}

	latency := time.Now().Sub(startTime).Milliseconds()
	l.Info("transfer transaction sent",
		slog.String("tx_hash", signedTx.Hash().Hex()),
		slog.Float64("latency", float64(latency)*0.001),
	)

	return signedTx.Hash().Hex(), nil
}
