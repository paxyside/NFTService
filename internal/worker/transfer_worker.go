package worker

import (
	"context"
	"errors"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rabbitmq/amqp091-go"
	"log/slog"
	"strings"
	"time"
)

func (w *Worker) TransferStatusUpdater() error {
	l := slog.Default()
	msgs, err := w.mq.Consume(w.tokenQueue.Name)
	if err != nil {
		return errors.New("failed to consume transfer update message")
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

			var txStatus string
			if receipt.Status == 1 {
				txStatus = "success"
			} else {
				txStatus = "failed"
			}

			if err := w.transferRepo.UpdateStatus(txStatus, txHash); err != nil {
				l.Error("failed to update transfer status", slog.Any("error", err))
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
