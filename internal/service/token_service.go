package service

import (
	"nft_service/infrastructure/utils"
	"nft_service/internal/contract"
	"nft_service/internal/domain"
)

type TokenService struct {
	repo     domain.TokenRepository
	contract *contract.NFTContract
}

func NewTokenService(repo domain.TokenRepository, contract *contract.NFTContract) *TokenService {
	return &TokenService{repo: repo, contract: contract}
}

func (t *TokenService) CreateToken(token *domain.Token) (*domain.Token, error) {

	var err error

	token.UniqueHash, err = utils.GenerateUniqueHash()
	if err != nil {
		return nil, err
	}

	if err := token.ValidateToCreate(); err != nil {
		return nil, err
	}

	token.TxHash, err = t.contract.Mint(token.Owner, token.UniqueHash, token.MediaUrl)
	if err != nil {
		return nil, err
	}

	if err := t.repo.CreateToken(token); err != nil {
		return nil, err
	}

	return token, nil
}