package domain

import (
	"testing"
)

func TestToken_ValidateToCreate(t *testing.T) {
	tests := []struct {
		name    string
		token   Token
		wantErr bool
	}{
		{
			name:    "Valid Token",
			token:   Token{MediaUrl: "https://example.com", Owner: "0x1234567890abcdef1234567890abcdef12345678"},
			wantErr: false,
		},
		{
			name:    "Empty MediaUrl",
			token:   Token{MediaUrl: "", Owner: "0x1234567890abcdef1234567890abcdef12345678"},
			wantErr: true,
		},
		{
			name:    "Invalid MediaUrl (not http/https)",
			token:   Token{MediaUrl: "ftp://example.com", Owner: "0x1234567890abcdef1234567890abcdef12345678"},
			wantErr: true,
		},
		{
			name:    "MediaUrl too long",
			token:   Token{MediaUrl: "https://" + string(make([]byte, 2049)), Owner: "0x1234567890abcdef1234567890abcdef12345678"},
			wantErr: true,
		},
		{
			name:    "Invalid Owner Address",
			token:   Token{MediaUrl: "https://example.com", Owner: "invalid_address"},
			wantErr: true,
		},
		{
			name:    "Valid Token with short Owner",
			token:   Token{MediaUrl: "https://example.com", Owner: "0x12345"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.token.ValidateToCreate()
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateToCreate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
