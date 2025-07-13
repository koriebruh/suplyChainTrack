package dto

import (
	"github.com/google/uuid"
	"github.com/koriebruh/suplyChainTrack/internal/domain"
	"time"
)

type CreateSupplyChainEventRequest struct {
	ProductID      *uuid.UUID   `json:"product_id"`
	StakeholderID  *uuid.UUID   `json:"stakeholder_id"`
	EventType      string       `json:"event_type" validate:"required,oneof=manufactured shipped received sold"`
	Location       *string      `json:"location"`
	Timestamp      time.Time    `json:"timestamp" validate:"required"`
	Metadata       domain.JSONB `json:"metadata"`
	BlockchainHash *string      `json:"blockchain_hash"`
}
