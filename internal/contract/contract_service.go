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
	"math/big"
	"nft_service/infrastructure/config"
	"strings"
)

type NFTContract struct {
	client      *ethclient.Client
	cfg         *config.Config
	contractABI string
}

func NewNFTContract(url string, cfg *config.Config, contractABI string) (*NFTContract, error) {
	client, err := ethclient.Dial(url)
	if err != nil {
		return nil, err
	}
	return &NFTContract{
		client:      client,
		cfg:         cfg,
		contractABI: contractABI,
	}, nil
}

func (m *NFTContract) Mint(owner string, uniqueHash string, mediaURL string) (string, error) {
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

	return signedTx.Hash().Hex(), nil
}

func (m *NFTContract) TotalSupply() (*big.Int, error) {

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
