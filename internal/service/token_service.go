package service

import (
	"encoding/json"
	"github.com/rabbitmq/amqp091-go"
	"math/big"
	"nft_service/infrastructure/rabbit"
	"nft_service/infrastructure/utils"
	"nft_service/internal/contract"
	"nft_service/internal/domain"
)

type TokenService struct {
	repo      domain.TokenRepository
	contract  contract.NFTService
	mq        *rabbit.RabbitMQ
	queueName amqp091.Queue
}

func NewTokenService(repo domain.TokenRepository, contract contract.NFTService, mq *rabbit.RabbitMQ, queueName amqp091.Queue) *TokenService {
	return &TokenService{repo: repo, contract: contract, mq: mq, queueName: queueName}
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

	token, err = t.contract.Mint(token)
	if err != nil {
		return nil, err
	}

	if err := t.repo.CreateToken(token); err != nil {
		return nil, err
	}

	queueBody, err := json.Marshal(token.TxHash)
	if err != nil {
		return nil, err
	}

	if err := t.mq.Publish(t.queueName.Name, queueBody); err != nil {
		return nil, err
	}

	return token, nil
}

func (t *TokenService) ListTokens(limit, offset int) ([]*domain.Token, error) {
	return t.repo.ListTokens(limit, offset)
}

func (t *TokenService) TotalSupply() (*big.Int, error) {
	return t.contract.TotalSupply()
}

func (t *TokenService) ExactTotalSupply() (*big.Int, error) {
	return t.contract.ExactTotalSupply()
}
