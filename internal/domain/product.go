package domain

import (
	"github.com/google/uuid"
	"time"
)

type Product struct {
	ID             uuid.UUID  `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	SKU            string     `json:"sku" gorm:"type:varchar(100);uniqueIndex;not null"`
	Name           string     `json:"name" gorm:"type:varchar(255);not null"`
	Description    *string    `json:"description" gorm:"type:text"`
	Category       *string    `json:"category" gorm:"type:varchar(100)"`
	ManufacturerID *uuid.UUID `json:"manufacturer_id" gorm:"type:uuid;index"`
	Metadata       JSONB      `json:"metadata" gorm:"type:jsonb"`
	CreatedAt      time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt      time.Time  `json:"updated_at" gorm:"autoUpdateTime"`

	// Relationships
	Manufacturer *Stakeholder `json:"manufacturer,omitempty" gorm:"foreignKey:ManufacturerID;constraint:OnDelete:SET NULL"`
}
