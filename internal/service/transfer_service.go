package service

import (
	"encoding/json"
	"github.com/rabbitmq/amqp091-go"
	"nft_service/infrastructure/rabbit"
	"nft_service/internal/contract"
	"nft_service/internal/domain"
)

type TransferService struct {
	repo      domain.TransferRepository
	contract  contract.NFTService
	mq        *rabbit.RabbitMQ
	queueName amqp091.Queue
}

func NewTransferService(repo domain.TransferRepository, contract contract.NFTService, mq *rabbit.RabbitMQ, queueName amqp091.Queue) *TransferService {
	return &TransferService{repo: repo, contract: contract, mq: mq, queueName: queueName}
}

func (s *TransferService) CreateTransfer(transfer *domain.Transfer) (*domain.Transfer, error) {
	var err error

	transfer, err = s.contract.TransferToken(transfer)
	if err != nil {
		return nil, err
	}

	err = s.repo.Create(transfer)
	if err != nil {
		return nil, err
	}

	queueBody, err := json.Marshal(transfer.TxHash)
	if err != nil {
		return nil, err
	}

	if err := s.mq.Publish(s.queueName.Name, queueBody); err != nil {
		return nil, err
	}

	return transfer, nil
}

func (s *TransferService) ListTransfer(limit, offset int) ([]domain.Transfer, error) {
	return s.repo.List(limit, offset)
}
