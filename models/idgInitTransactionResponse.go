package models

type IdgInitTransactionResponse struct {
	Status      string      `json:"status"`
	Message     interface{} `json:"message"`
	ErrorCode   interface{} `json:"error_code"`
	StatusCode  interface{} `json:"status_code"`
	IsRetryable bool        `json:"is_retryable"`
	Token       string      `json:"token"`
}
