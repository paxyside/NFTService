package persistence

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"nft_service/internal/domain"
	"strings"
)

type TokenRepo struct {
	db *pgxpool.Pool
}

func NewTokenRepo(db *pgxpool.Pool) *TokenRepo {
	return &TokenRepo{db: db}
}

func (t TokenRepo) CreateToken(token *domain.Token) error {

	var tokenId sql.NullString

	query := `INSERT INTO nfts (unique_hash, tx_hash, media_url, owner)
			  VALUES ($1, $2, $3, $4)
			  RETURNING id, unique_hash, tx_hash, media_url, owner, token_id, created_at`

	err := t.db.QueryRow(context.Background(), query, token.UniqueHash, token.TxHash, token.MediaUrl, token.Owner).Scan(
		&token.ID,
		&token.UniqueHash,
		&token.TxHash,
		&token.MediaUrl,
		&token.Owner,
		&tokenId,
		&token.CreatedAt,
	)

	if err != nil {
		if strings.Contains(err.Error(), "duplicate") {
			return errors.New("token already exists")
		}
		return fmt.Errorf("failed to create token: %w", err)
	}

	return nil
}

func (t TokenRepo) UpdateTokenIDWithTransaction(tokenID, txHash string) error {
	tx, err := t.db.Begin(context.Background())
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback(context.Background())
		}
	}()

	query := `UPDATE nfts SET token_id = $1 WHERE tx_hash = $2`
	row, err := tx.Exec(context.Background(), query, tokenID, txHash)
	if err != nil {
		return fmt.Errorf("failed to update token id: %w", err)
	}

	if row.RowsAffected() == 0 {
		return errors.New("token with this tx_hash does not exist")
	}

	if err = tx.Commit(context.Background()); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (t TokenRepo) ListTokens(limit, offset int) ([]*domain.Token, error) {

	var tokens []*domain.Token

	query := `SELECT * FROM nfts LIMIT $1 OFFSET $2`

	rows, err := t.db.Query(context.TODO(), query, limit, offset)
	defer rows.Close()

	if err != nil {
		switch {
		case strings.Contains(err.Error(), "no rows"):
			return tokens, nil
		default:
			return nil, errors.New("token receipt error " + err.Error())
		}
	}

	for rows.Next() {
		token := &domain.Token{}

		err := rows.Scan(
			&token.ID,
			&token.UniqueHash,
			&token.TxHash,
			&token.MediaUrl,
			&token.Owner,
			&token.CreatedAt,
			&token.TokenID,
		)
		if err != nil {
			return nil, errors.New("scan error " + err.Error())
		}

		tokens = append(tokens, token)
	}

	if err := rows.Err(); err != nil {
		return nil, errors.New("rows scan error " + err.Error())
	}

	return tokens, err
}
