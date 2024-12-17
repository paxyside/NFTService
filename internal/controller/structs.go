package controller

type CreateTokenRequest struct {
	Owner    string `json:"owner"`
	MediaUrl string `json:"media_url"`
}

type SupplyResponse struct {
	TotalSupply string `json:"total_supply"`
}

type ErrorResponse struct {
	RequestID string `json:"request_id"`
	Error     string `json:"error"`
}
type CreateTransferRequest struct {
	From    string `json:"from_address"`
	To      string `json:"to_address"`
	TokenId string `json:"token_id"`
}
