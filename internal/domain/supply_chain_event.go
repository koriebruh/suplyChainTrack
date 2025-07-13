package domain

import (
	"github.com/google/uuid"
	"time"
)

// EventType constants
const (
	EventTypeManufactured = "manufactured"
	EventTypeShipped      = "shipped"
	EventTypeReceived     = "received"
	EventTypeSold         = "sold"
)

func IsValidEventType(t string) bool {
	switch t {
	case EventTypeManufactured,
		EventTypeShipped,
		EventTypeReceived,
		EventTypeSold:
		return true
	default:
		return false
	}
}

type SupplyChainEvent struct {
	ID             uuid.UUID  `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	ProductID      *uuid.UUID `json:"product_id" gorm:"type:uuid;index"`
	StakeholderID  *uuid.UUID `json:"stakeholder_id" gorm:"type:uuid;index"`
	EventType      string     `json:"event_type" gorm:"type:varchar(50);not null"`
	Location       *string    `json:"location" gorm:"type:varchar(255)"`
	Timestamp      time.Time  `json:"timestamp" gorm:"not null"`
	Metadata       JSONB      `json:"metadata" gorm:"type:jsonb"`
	BlockchainHash *string    `json:"blockchain_hash" gorm:"type:varchar(66)"`
	IsVerified     bool       `json:"is_verified" gorm:"default:false"`
	CreatedAt      time.Time  `json:"created_at" gorm:"autoCreateTime"`

	// Relationships
	Product     *Product     `json:"product,omitempty" gorm:"foreignKey:ProductID;constraint:OnDelete:CASCADE"`
	Stakeholder *Stakeholder `json:"stakeholder,omitempty" gorm:"foreignKey:StakeholderID;constraint:OnDelete:CASCADE"`
}
