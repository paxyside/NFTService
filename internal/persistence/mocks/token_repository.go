package mocks

import (
	"errors"
	"github.com/stretchr/testify/mock"
	"nft_service/internal/domain"
	"sort"
	"sync"
)

type MockTokenRepository struct {
	mock.Mock
	tokens map[int]*domain.Token
	mu     sync.RWMutex
}

func NewMockTokenRepository() *MockTokenRepository {
	return &MockTokenRepository{
		tokens: make(map[int]*domain.Token),
	}
}

func (m *MockTokenRepository) CreateToken(token *domain.Token) error {
	if token == nil {
		return errors.New("token cannot be nil")
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	for _, existingToken := range m.tokens {
		if existingToken.UniqueHash == token.UniqueHash {
			return errors.New("token already exists")
		}
	}

	if token.UniqueHash == "simulate_error" {
		return errors.New("database error: simulated failure")
	}

	m.tokens[token.ID] = token
	return nil
}

func (m *MockTokenRepository) ListTokens(limit, offset int) ([]*domain.Token, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if limit <= 0 || offset < 0 {
		return nil, errors.New("invalid limit or offset")
	}

	sortedTokens := make([]*domain.Token, 0, len(m.tokens))
	for _, token := range m.tokens {
		sortedTokens = append(sortedTokens, token)
	}
	sort.Slice(sortedTokens, func(i, j int) bool {
		return sortedTokens[i].ID < sortedTokens[j].ID
	})

	start := offset
	if start >= len(sortedTokens) {
		return []*domain.Token{}, nil
	}
	end := start + limit
	if end > len(sortedTokens) {
		end = len(sortedTokens)
	}

	return sortedTokens[start:end], nil
}
