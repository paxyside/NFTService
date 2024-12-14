package persistence

import (
	"context"
	"errors"
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

	query := `INSERT INTO nfts (unique_hash, tx_hash, media_url, owner)
			  VALUES ($1, $2, $3, $4)
			  RETURNING id, unique_hash, tx_hash, media_url, owner, created_at`

	err := t.db.QueryRow(context.Background(), query, token.UniqueHash, token.TxHash, token.MediaUrl, token.Owner).Scan(
		&token.ID,
		&token.UniqueHash,
		&token.TxHash,
		&token.MediaUrl,
		&token.Owner,
		&token.CreatedAt,
	)

	if err != nil {
		if strings.Contains(err.Error(), "duplicate") {
			return errors.New("token already exists")
		}
		return errors.New("failed to create token: " + err.Error())
	}

	return nil
}
