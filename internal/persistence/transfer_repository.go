package persistence

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"nft_service/internal/domain"
	"strings"
)

type TransferRepo struct {
	db *pgxpool.Pool
}

func NewTransferRepo(db *pgxpool.Pool) *TransferRepo {
	return &TransferRepo{
		db: db,
	}
}

func (t TransferRepo) Create(transfer *domain.Transfer) error {

	query := `INSERT INTO transfers (from_address, to_address, token_id, tx_hash, status)
			  VALUES ($1, $2, $3, $4, $5)
              RETURNING id, from_address, to_address, token_id, tx_hash, status, created_at, updated_at`

	err := t.db.QueryRow(context.Background(), query, transfer.FromAddress, transfer.ToAddress, transfer.TokenID, transfer.TxHash, transfer.Status).Scan(
		&transfer.ID,
		&transfer.FromAddress,
		&transfer.ToAddress,
		&transfer.TokenID,
		&transfer.TxHash,
		&transfer.Status,
		&transfer.CreatedAt,
		&transfer.UpdatedAt,
	)

	if err != nil {
		if strings.Contains(err.Error(), "duplicate") {
			return errors.New("transfer already exists")
		}
		return fmt.Errorf("failed to create transfer: %w", err)
	}

	return nil
}

func (t TransferRepo) UpdateStatus(status, txHash string) error {
	tx, err := t.db.Begin(context.Background())
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback(context.Background())
		}
	}()

	query := `UPDATE transfers SET status = $1 WHERE tx_hash = $2`
	row, err := tx.Exec(context.Background(), query, status, txHash)
	if err != nil {
		return fmt.Errorf("failed to update transfer status: %w", err)
	}

	if row.RowsAffected() == 0 {
		return errors.New("transfer with this tx_hash does not exist")
	}

	if err = tx.Commit(context.Background()); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (t TransferRepo) List(limit, offset int) ([]domain.Transfer, error) {
	var transfers []domain.Transfer

	query := `SELECT * FROM transfers LIMIT $1 OFFSET $2`
	rows, err := t.db.Query(context.Background(), query, limit, offset)

	defer rows.Close()

	if err != nil {
		switch {
		case strings.Contains(err.Error(), "no rows"):
			return transfers, nil
		default:
			return nil, errors.New("token receipt error " + err.Error())
		}
	}

	for rows.Next() {
		transfer := domain.Transfer{}
		err := rows.Scan(
			&transfer.ID,
			&transfer.FromAddress,
			&transfer.ToAddress,
			&transfer.TokenID,
			&transfer.TxHash,
			&transfer.Status,
			&transfer.CreatedAt,
			&transfer.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan transfer row: %w", err)
		}
		transfers = append(transfers, transfer)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate transfers: %w", err)
	}

	return transfers, nil
}
