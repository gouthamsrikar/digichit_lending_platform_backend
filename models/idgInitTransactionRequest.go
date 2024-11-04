package models

type IdgInitTransactionRequest struct {
	GatewayInstanceID string        `json:"gateway_instance_id"`
	Keys              []interface{} `json:"keys"`
	IdempotencyID     string        `json:"idempotency_id"`
	ClientID          string        `json:"client_id"`
	ClientSecret      string        `json:"client_secret"`
	TransactionInput  struct {
		Mobile string `json:"mobile"`
	} `json:"transaction_input"`
}
