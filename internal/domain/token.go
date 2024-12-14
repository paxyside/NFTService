package domain

import (
	"errors"
	"net/url"
	"regexp"
	"time"
)

const (
	ethereumOwnerAddressExpression = `^0x[a-fA-F0-9]{40}$`
)

type TokenRepository interface {
	CreateToken(token *Token) error
	ListTokens(limit, offset int) ([]*Token, error)
}

type Token struct {
	ID         int       `json:"id,omitempty"`
	UniqueHash string    `json:"unique_hash,omitempty"`
	TxHash     string    `json:"tx_hash,omitempty"`
	MediaUrl   string    `json:"media_url" binding:"required"`
	Owner      string    `json:"owner" binding:"required"`
	CreatedAt  time.Time `json:"created_at,omitempty"`
}

func (t *Token) ValidateToCreate() error {

	var (
		rgx *regexp.Regexp
		err error
	)

	if t.MediaUrl == "" || len(t.MediaUrl) > 2048 {
		return errors.New("invalid media_url, must be non-empty and less than 2048 characters")
	}

	u, err := url.Parse(t.MediaUrl)
	if err != nil || (u.Scheme != "http" && u.Scheme != "https") {
		return errors.New("invalid media_url, must be a valid http or https URL")
	}

	rgx, err = regexp.Compile(ethereumOwnerAddressExpression)

	if err != nil {
		return errors.New("failed to compile Ethereum address regex: " + err.Error())
	}

	if !rgx.MatchString(t.Owner) {
		return errors.New("invalid owner address " + t.Owner)
	}

	return nil
}
