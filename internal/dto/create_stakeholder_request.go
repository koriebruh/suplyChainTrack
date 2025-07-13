package dto

type CreateStakeholderRequest struct {
	Name          string  `json:"name" validate:"required"`
	Type          string  `json:"type" validate:"required,oneof=manufacturer distributor retailer"`
	WalletAddress *string `json:"wallet_address"`
	Email         string  `json:"email" validate:"required,email"`
	Phone         *string `json:"phone"`
	Address       *string `json:"address"`
}
