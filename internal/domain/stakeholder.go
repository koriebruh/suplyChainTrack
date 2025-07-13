package domain

import (
	"github.com/google/uuid"
	"time"
)

// StakeholderType constants
const (
	StakeholderTypeManufacturer = "manufacturer"
	StakeholderTypeDistributor  = "distributor"
	StakeholderTypeRetailer     = "retailer"
)

func IsValidStakeholderType(t string) bool {
	switch t {
	case StakeholderTypeManufacturer,
		StakeholderTypeDistributor,
		StakeholderTypeRetailer:
		return true
	default:
		return false
	}
}

type Stakeholder struct {
	ID            uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Name          string    `json:"name" gorm:"type:varchar(255);not null"`
	Type          string    `json:"type" gorm:"type:varchar(50);not null"` // 'manufacturer', 'distributor', 'retailer'
	WalletAddress *string   `json:"wallet_address" gorm:"type:varchar(42);uniqueIndex"`
	Email         string    `json:"email" gorm:"type:varchar(255);uniqueIndex;not null"`
	Phone         *string   `json:"phone" gorm:"type:varchar(20)"`
	Address       *string   `json:"address" gorm:"type:text"`
	IsVerified    bool      `json:"is_verified" gorm:"default:false"`
	CreatedAt     time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt     time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}
