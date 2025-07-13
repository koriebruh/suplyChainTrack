package dto

import (
	"github.com/google/uuid"
	"github.com/koriebruh/suplyChainTrack/internal/domain"
	"time"
)

// SupplyChainTrace represents a complete trace of product through supply chain
type SupplyChainTrace struct {
	Product *domain.Product            `json:"product"`
	Events  []*domain.SupplyChainEvent `json:"events"`
}

// StakeholderStats represents statistics for a stakeholder
type StakeholderStats struct {
	StakeholderID  uuid.UUID `json:"stakeholder_id"`
	TotalProducts  int       `json:"total_products"`
	TotalEvents    int       `json:"total_events"`
	VerifiedEvents int       `json:"verified_events"`
	PendingEvents  int       `json:"pending_events"`
	LastActivity   time.Time `json:"last_activity"`
}

// ProductStats represents statistics for a product
type ProductStats struct {
	ProductID       uuid.UUID  `json:"product_id"`
	TotalEvents     int        `json:"total_events"`
	VerifiedEvents  int        `json:"verified_events"`
	CurrentLocation *string    `json:"current_location"`
	LastStakeholder *uuid.UUID `json:"last_stakeholder"`
	LastActivity    time.Time  `json:"last_activity"`
}

// Filter untuk query
type StakeholderFilter struct {
	Type       *string `json:"type"`
	IsVerified *bool   `json:"is_verified"`
	Email      *string `json:"email"`
	Limit      int     `json:"limit"`
	Offset     int     `json:"offset"`
}

type ProductFilter struct {
	Category       *string    `json:"category"`
	ManufacturerID *uuid.UUID `json:"manufacturer_id"`
	SKU            *string    `json:"sku"`
	Name           *string    `json:"name"`
	Limit          int        `json:"limit"`
	Offset         int        `json:"offset"`
}

type SupplyChainEventFilter struct {
	ProductID     *uuid.UUID `json:"product_id"`
	StakeholderID *uuid.UUID `json:"stakeholder_id"`
	EventType     *string    `json:"event_type"`
	Location      *string    `json:"location"`
	IsVerified    *bool      `json:"is_verified"`
	FromDate      *time.Time `json:"from_date"`
	ToDate        *time.Time `json:"to_date"`
	Limit         int        `json:"limit"`
	Offset        int        `json:"offset"`
}

type BlockchainTransactionFilter struct {
	EventID *uuid.UUID `json:"event_id"`
	Status  *string    `json:"status"`
	Limit   int        `json:"limit"`
	Offset  int        `json:"offset"`
}

// Response wrapper
type PaginatedResponse struct {
	Data    interface{} `json:"data"`
	Total   int         `json:"total"`
	Limit   int         `json:"limit"`
	Offset  int         `json:"offset"`
	HasMore bool        `json:"has_more"`
}

// Error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
	Code    int    `json:"code"`
}
