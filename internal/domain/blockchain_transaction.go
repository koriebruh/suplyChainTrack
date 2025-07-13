package domain

import (
	"github.com/google/uuid"
	"time"
)

// TransactionStatus constants
const (
	TransactionStatusPending   = "pending"
	TransactionStatusConfirmed = "confirmed"
	TransactionStatusFailed    = "failed"
)

type BlockchainTransaction struct {
	ID              uuid.UUID  `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	EventID         *uuid.UUID `json:"event_id" gorm:"type:uuid;index"`
	TransactionHash string     `json:"transaction_hash" gorm:"type:varchar(66);uniqueIndex;not null"`
	BlockNumber     *int64     `json:"block_number" gorm:"type:bigint"`
	GasUsed         *int64     `json:"gas_used" gorm:"type:bigint"`
	Status          string     `json:"status" gorm:"type:varchar(20);default:'pending'"`
	CreatedAt       time.Time  `json:"created_at" gorm:"autoCreateTime"`

	// Relationships
	Event *SupplyChainEvent `json:"event,omitempty" gorm:"foreignKey:EventID;constraint:OnDelete:CASCADE"`
}
