package domain

import (
	"errors"
	"math/big"
	"regexp"
	"time"
)

type TransferRepository interface {
	Create(transfer *Transfer) error
	UpdateStatus(status, txHash string) error
	List(limit, offset int) ([]Transfer, error)
}

type Transfer struct {
	ID          int       `json:"id"`
	FromAddress string    `json:"from_address" binding:"required"`
	ToAddress   string    `json:"to_address" binding:"required"`
	TokenID     string    `json:"token_id" binding:"required"`
	TxHash      string    `json:"tx_hash"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (t *Transfer) ValidateToCreate() error {

	var (
		rgx *regexp.Regexp
		err error
	)

	rgx, err = regexp.Compile(ethereumAddressExpression)

	if err != nil {
		return errors.New("failed to compile Ethereum address regex: " + err.Error())
	}

	if !rgx.MatchString(t.FromAddress) {
		return errors.New("invalid from address " + t.FromAddress)
	}

	if !rgx.MatchString(t.ToAddress) {
		return errors.New("invalid to address " + t.ToAddress)
	}

	if t.TokenID == "" {
		return errors.New("invalid token id")
	}

	_, ok := new(big.Int).SetString(t.TokenID, 10)
	if !ok {
		return errors.New("invalid token id")
	}

	return nil
}
