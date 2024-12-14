package persistence

import (
	"github.com/stretchr/testify/assert"
	"nft_service/internal/domain"
	"nft_service/internal/persistence/mocks"
	"testing"
)

func TestMockTokenRepository_CreateToken(t *testing.T) {
	repo := mocks.NewMockTokenRepository()

	token := &domain.Token{
		ID:         1,
		UniqueHash: "unique_hash_1",
		TxHash:     "tx_hash_1",
		MediaUrl:   "https://example.com/media1",
		Owner:      "owner_1",
	}

	err := repo.CreateToken(token)
	assert.NoError(t, err)

	tokenDuplicate := &domain.Token{
		ID:         2,
		UniqueHash: "unique_hash_1",
		TxHash:     "tx_hash_2",
		MediaUrl:   "https://example.com/media2",
		Owner:      "owner_2",
	}

	err = repo.CreateToken(tokenDuplicate)
	assert.EqualError(t, err, "token already exists")

	errorToken := &domain.Token{
		ID:         3,
		UniqueHash: "simulate_error",
		TxHash:     "tx_hash_3",
		MediaUrl:   "https://example.com/media3",
		Owner:      "owner_3",
	}

	err = repo.CreateToken(errorToken)
	assert.EqualError(t, err, "database error: simulated failure")

	err = repo.CreateToken(nil)
	assert.EqualError(t, err, "token cannot be nil")
}

func TestMockTokenRepository_ListTokens(t *testing.T) {
	repo := mocks.NewMockTokenRepository()

	repo.CreateToken(&domain.Token{
		ID:         1,
		UniqueHash: "unique_hash_1",
		TxHash:     "tx_hash_1",
		MediaUrl:   "https://example.com/media1",
		Owner:      "owner_1",
	})
	repo.CreateToken(&domain.Token{
		ID:         2,
		UniqueHash: "unique_hash_2",
		TxHash:     "tx_hash_2",
		MediaUrl:   "https://example.com/media2",
		Owner:      "owner_2",
	})
	repo.CreateToken(&domain.Token{
		ID:         3,
		UniqueHash: "unique_hash_3",
		TxHash:     "tx_hash_3",
		MediaUrl:   "https://example.com/media3",
		Owner:      "owner_3",
	})

	tokens, err := repo.ListTokens(10, 0)
	assert.NoError(t, err)
	assert.Len(t, tokens, 3)
	assert.Equal(t, "unique_hash_1", tokens[0].UniqueHash)
	assert.Equal(t, "unique_hash_3", tokens[2].UniqueHash)

	tokens, err = repo.ListTokens(2, 1)
	assert.NoError(t, err)
	assert.Len(t, tokens, 2)
	assert.Equal(t, "unique_hash_2", tokens[0].UniqueHash)

	tokens, err = repo.ListTokens(10, 10)
	assert.NoError(t, err)
	assert.Empty(t, tokens)

	tokens, err = repo.ListTokens(-1, 0)
	assert.Nil(t, tokens)
	assert.EqualError(t, err, "invalid limit or offset")

	tokens, err = repo.ListTokens(10, -1)
	assert.Nil(t, tokens)
	assert.EqualError(t, err, "invalid limit or offset")
}
