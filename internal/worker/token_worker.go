package worker

import (
	"context"
	"errors"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rabbitmq/amqp091-go"
	"log/slog"
	"math/big"
	"strings"
	"time"
)

func (w *Worker) TokenUpdater() error {
	l := slog.Default()
	msgs, err := w.mq.Consume(w.tokenQueue.Name)
	if err != nil {
		return errors.New("failed to consume token update message")
	}

	maxWorkers := 10
	semaphore := make(chan struct{}, maxWorkers)

	for msg := range msgs {
		semaphore <- struct{}{}

		go func(msg amqp091.Delivery) {
			defer func() { <-semaphore }()

			txHash := strings.Trim(string(msg.Body), "\"")

			ctx, cancel := context.WithTimeout(context.Background(), 55*time.Second)
			defer cancel()

			receipt, err := w.client.TransactionReceipt(ctx, common.HexToHash(txHash))
			if err != nil {
				select {
				case <-ctx.Done():
					return
				default:
					time.Sleep(5 * time.Second)
					msg.Nack(false, true)
					return
				}
			}

			var tokenID string
			for _, log := range receipt.Logs {
				if log.Topics[0].Hex() == w.parsedABI.Events["Transfer"].ID.Hex() {
					var transferEvent struct {
						From    common.Address
						To      common.Address
						TokenID *big.Int
					}
					if err := w.parsedABI.UnpackIntoInterface(&transferEvent, "Transfer", log.Data); err != nil {
						l.Error("failed to unpack transfer event", slog.Any("error", err))
						return
					}
					tokenID = new(big.Int).SetBytes(log.Topics[3].Bytes()).String()
					break
				}
			}

			if err := w.repo.UpdateTokenID(tokenID, txHash); err != nil {
				l.Error("failed to update token", slog.Any("error", err))
				msg.Nack(false, true)
				return
			}

			if err := msg.Ack(false); err != nil {
				l.Error("failed to ack message", slog.Any("error", err))
			}
		}(msg)
	}

	return nil
}
