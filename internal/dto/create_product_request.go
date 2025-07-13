package dto

import (
	"github.com/google/uuid"
	"github.com/koriebruh/suplyChainTrack/internal/domain"
)

type CreateProductRequest struct {
	SKU            string       `json:"sku" validate:"required"`
	Name           string       `json:"name" validate:"required"`
	Description    *string      `json:"description"`
	Category       *string      `json:"category"`
	ManufacturerID *uuid.UUID   `json:"manufacturer_id"`
	Metadata       domain.JSONB `json:"metadata"`
}
