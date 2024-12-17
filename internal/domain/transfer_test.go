package domain

import (
	"testing"
)

func TestValidateToCreate(t *testing.T) {
	tests := []struct {
		transfer  Transfer
		expectErr bool
	}{
		{
			transfer: Transfer{
				FromAddress: "0xC92f65c05ccdeF650fe1fdeC0221E5f993ea8956",
				ToAddress:   "0xC92f65c05ccdeF650fe1fdeC0221E5f993ea8956",
				TokenID:     "123",
			},
			expectErr: false,
		},
		{
			transfer: Transfer{
				FromAddress: "invalid_address",
				ToAddress:   "0xC92f65c05ccdeF650fe1fdeC0221E5f993ea8956",
				TokenID:     "123",
			},
			expectErr: true,
		},
		{
			transfer: Transfer{
				FromAddress: "0xC92f65c05ccdeF650fe1fdeC0221E5f993ea8956",
				ToAddress:   "invalid_address",
				TokenID:     "123",
			},
			expectErr: true,
		},
		{
			transfer: Transfer{
				FromAddress: "0xC92f65c05ccdeF650fe1fdeC0221E5f993ea8956",
				ToAddress:   "0xC92f65c05ccdeF650fe1fdeC0221E5f993ea8956",
				TokenID:     "",
			},
			expectErr: true,
		},
		{
			transfer: Transfer{
				FromAddress: "0xC92f65c05ccdeF650fe1fdeC0221E5f993ea8956",
				ToAddress:   "0xC92f65c05ccdeF650fe1fdeC0221E5f993ea8956",
				TokenID:     "abc",
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.transfer.TokenID, func(t *testing.T) {
			err := tt.transfer.ValidateToCreate()
			if (err != nil) != tt.expectErr {
				t.Errorf("expected error: %v, got: %v", tt.expectErr, err)
			}
		})
	}
}
