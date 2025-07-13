package dto

type UpdateStakeholderRequest struct {
	Name          *string `json:"name"`
	WalletAddress *string `json:"wallet_address"`
	Email         *string `json:"email" validate:"omitempty,email"`
	Phone         *string `json:"phone"`
	Address       *string `json:"address"`
	IsVerified    *bool   `json:"is_verified"`
}
