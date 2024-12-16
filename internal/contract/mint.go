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
	"nft_service/internal/domain"
	"strings"
	"time"
)

func (m *NFTContract) Mint(token *domain.Token) (*domain.Token, error) {

	var (
		l         = slog.Default()
		startTime = time.Now()
	)

	parsedAbi, err := abi.JSON(strings.NewReader(m.contractABI))
	if err != nil {
		return nil, fmt.Errorf("failed to parse contract ABI: %w", err)
	}

	txData, err := parsedAbi.Pack("mint", common.HexToAddress(token.Owner), token.UniqueHash, token.MediaUrl)
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
		return nil, fmt.Errorf("failed to convert private key: %w", err)
	}

	signedTx, err := types.SignTx(unsignedTx, types.LatestSignerForChainID(big.NewInt(m.cfg.ChainID)), privateKeyBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to sign transaction: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err = m.client.SendTransaction(ctx, signedTx)
	if err != nil {
		return nil, fmt.Errorf("failed to send transaction: %w", err)
	}

	latency1 := time.Now().Sub(startTime).Milliseconds()
	l.Info("mint transaction sent to contract", slog.Float64("latency", float64(latency1)*0.001))

	var receipt *types.Receipt
	for {
		receipt, err = m.client.TransactionReceipt(ctx, signedTx.Hash())
		if err == nil {
			break
		}
		if err.Error() != "not found" {
			return nil, fmt.Errorf("failed to get transaction receipt: %w", err)
		}
		time.Sleep(2 * time.Second)
	}

	for _, log := range receipt.Logs {
		if log.Topics[0].Hex() == parsedAbi.Events["Transfer"].ID.Hex() {
			var transferEvent struct {
				From    common.Address
				To      common.Address
				TokenID *big.Int
			}
			if err := parsedAbi.UnpackIntoInterface(&transferEvent, "Transfer", log.Data); err != nil {
				return nil, fmt.Errorf("failed to unpack event data: %w", err)
			}
			transferEvent.TokenID = new(big.Int).SetBytes(log.Topics[3].Bytes())
			token.TokenID = transferEvent.TokenID
			break
		}
	}

	token.TxHash = signedTx.Hash().Hex()

	latency2 := time.Now().Sub(startTime).Milliseconds()
	l.Info("mint transaction submitted with tokenID", slog.Float64("latency", float64(latency2)*0.001))

	return token, nil
}
