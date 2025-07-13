package dto

import (
	"github.com/google/uuid"
	"github.com/koriebruh/suplyChainTrack/internal/domain"
)

type UpdateProductRequest struct {
	Name           *string      `json:"name"`
	Description    *string      `json:"description"`
	Category       *string      `json:"category"`
	ManufacturerID *uuid.UUID   `json:"manufacturer_id"`
	Metadata       domain.JSONB `json:"metadata"`
}
