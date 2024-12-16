package domain

import (
	"errors"
	"regexp"
)

type TransferRepository interface {
	Create(transfer *Transfer) error
	TransferList() ([]Transfer, error)
}

type Transfer struct {
	ID          int    `json:"id"`
	FromAddress string `json:"from_address"`
	ToAddress   string `json:"to_address"`
	TokenID     string `json:"token_id"`
	TxHash      string `json:"tx_hash"`
	Status      string `json:"status"`
	CreatedAt   int64  `json:"created_at"`
	UpdatedAt   int64  `json:"updated_at"`
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

	return nil
}
