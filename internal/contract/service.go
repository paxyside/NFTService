package contract

import (
	"github.com/ethereum/go-ethereum/ethclient"
	"math/big"
	"nft_service/infrastructure/config"
	"nft_service/internal/domain"
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
