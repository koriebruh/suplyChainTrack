package dto

import "github.com/google/uuid"

type CreateBlockchainTransactionRequest struct {
	EventID         *uuid.UUID `json:"event_id"`
	TransactionHash string     `json:"transaction_hash" validate:"required"`
	BlockNumber     *int64     `json:"block_number"`
	GasUsed         *int64     `json:"gas_used"`
	Status          string     `json:"status" validate:"required,oneof=pending confirmed failed"`
}
