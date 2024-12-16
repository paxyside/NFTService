package contract

import (
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/ethclient"
	"math/big"
	"nft_service/infrastructure/config"
	"nft_service/internal/domain"
	"strings"
	"sync"
)

type NFTService interface {
	Mint(token *domain.Token) (*domain.Token, error)
	TotalSupply() (*big.Int, error)
	ExactTotalSupply() (*big.Int, error)
	TransferToken(from string, to string, tokenId *big.Int) (string, error)
}

type NFTContract struct {
	client       *ethclient.Client
	cfg          *config.Config
	contractABI  string
	parsedABI    *abi.ABI
	cache        *big.Int
	cacheUpdated int64
	mu           sync.RWMutex
}

func NewNFTContract(url string, cfg *config.Config, contractABI string) (*NFTContract, error) {
	client, err := ethclient.Dial(url)
	if err != nil {
		return nil, err
	}

	parsedAbi, err := abi.JSON(strings.NewReader(contractABI))
	if err != nil {
		return nil, fmt.Errorf("failed to parse contract ABI: %w", err)
	}

	contract := &NFTContract{
		client:    client,
		cfg:       cfg,
		parsedABI: &parsedAbi,
	}

	return contract, nil
}
