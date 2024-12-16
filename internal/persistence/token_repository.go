package persistence

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"math/big"
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
	var tokenID int64
	if token.TokenID != nil {
		if !token.TokenID.IsInt64() {
			return fmt.Errorf("token_id value is too large for BIGINT")
		}
		tokenID = token.TokenID.Int64()
	}

	query := `INSERT INTO nfts (unique_hash, tx_hash, media_url, owner, token_id)
			  VALUES ($1, $2, $3, $4, $5)
			  RETURNING id, unique_hash, tx_hash, media_url, owner, token_id, created_at`

	var returnedTokenID int64
	err := t.db.QueryRow(context.Background(), query, token.UniqueHash, token.TxHash, token.MediaUrl, token.Owner, tokenID).Scan(
		&token.ID,
		&token.UniqueHash,
		&token.TxHash,
		&token.MediaUrl,
		&token.Owner,
		&returnedTokenID,
		&token.CreatedAt,
	)

	if err != nil {
		if strings.Contains(err.Error(), "duplicate") {
			return errors.New("token already exists")
		}
		return fmt.Errorf("failed to create token: %w", err)
	}

	token.TokenID = new(big.Int).SetInt64(returnedTokenID)

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
		var tokenID int64

		err := rows.Scan(
			&token.ID,
			&token.UniqueHash,
			&token.TxHash,
			&token.MediaUrl,
			&token.Owner,
			&token.CreatedAt,
			&tokenID,
		)
		if err != nil {
			return nil, errors.New("scan error " + err.Error())
		}

		token.TokenID = big.NewInt(tokenID)

		tokens = append(tokens, token)
	}

	if err := rows.Err(); err != nil {
		return nil, errors.New("rows scan error " + err.Error())
	}

	return tokens, err
}
