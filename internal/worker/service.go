package worker

import (
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/rabbitmq/amqp091-go"
	"nft_service/infrastructure/rabbit"
	"nft_service/internal/domain"
	"strings"
)

type WorkerUpdater interface {
	TokenUpdater() error
}

type Worker struct {
	client      *ethclient.Client
	mq          *rabbit.RabbitMQ
	tokenQueue  amqp091.Queue
	repo        domain.TokenRepository
	contractABI string
	parsedABI   *abi.ABI
}

func NewWorker(url string, mq *rabbit.RabbitMQ, tokenQueue amqp091.Queue, repo domain.TokenRepository, contractABI string) (*Worker, error) {
	client, err := ethclient.Dial(url)
	if err != nil {
		return nil, err
	}

	parsedAbi, err := abi.JSON(strings.NewReader(contractABI))
	if err != nil {
		return nil, fmt.Errorf("failed to parse contract ABI: %w", err)
	}

	return &Worker{
		client:     client,
		mq:         mq,
		tokenQueue: tokenQueue,
		repo:       repo,
		parsedABI:  &parsedAbi,
	}, nil
}
