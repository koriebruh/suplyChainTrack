package repository

import (
	"context"
	"github.com/google/uuid"
	"github.com/koriebruh/suplyChainTrack/internal/domain"
	"github.com/koriebruh/suplyChainTrack/internal/dto"
	"gorm.io/gorm"
)

type Repositories struct {
	Stakeholder           StakeholderRepository
	Product               ProductRepository
	SupplyChainEvent      SupplyChainEventRepository
	BlockchainTransaction BlockchainTransactionRepository
}

func NewRepositories(db *gorm.DB) *Repositories {
	return &Repositories{
		Stakeholder:           NewStakeholderRepository(db),
		Product:               NewProductRepository(db),
		SupplyChainEvent:      NewSupplyChainEventRepository(db),
		BlockchainTransaction: NewBlockchainTransactionRepository(db),
	}
}

type StakeholderRepository interface {
	Create(ctx context.Context, stakeholder *domain.Stakeholder) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Stakeholder, error)
	GetByEmail(ctx context.Context, email string) (*domain.Stakeholder, error)
	GetByWalletAddress(ctx context.Context, address string) (*domain.Stakeholder, error)
	Update(ctx context.Context, id uuid.UUID, updates map[string]interface{}) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, filter *dto.StakeholderFilter) ([]*domain.Stakeholder, int64, error)
	GetStats(ctx context.Context, id uuid.UUID) (*dto.StakeholderStats, error)
}

type ProductRepository interface {
	Create(ctx context.Context, product *domain.Product) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Product, error)
	GetBySKU(ctx context.Context, sku string) (*domain.Product, error)
	Update(ctx context.Context, id uuid.UUID, updates map[string]interface{}) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, filter *dto.ProductFilter) ([]*domain.Product, int64, error)
	GetStats(ctx context.Context, id uuid.UUID) (*dto.ProductStats, error)
	GetByManufacturer(ctx context.Context, manufacturerID uuid.UUID) ([]*domain.Product, error)
}

type SupplyChainEventRepository interface {
	Create(ctx context.Context, event *domain.SupplyChainEvent) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.SupplyChainEvent, error)
	Update(ctx context.Context, id uuid.UUID, updates map[string]interface{}) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, filter *dto.SupplyChainEventFilter) ([]*domain.SupplyChainEvent, int64, error)
	GetByProduct(ctx context.Context, productID uuid.UUID) ([]*domain.SupplyChainEvent, error)
	GetByStakeholder(ctx context.Context, stakeholderID uuid.UUID) ([]*domain.SupplyChainEvent, error)
	GetTrace(ctx context.Context, productID uuid.UUID) (*dto.SupplyChainTrace, error)
	VerifyEvent(ctx context.Context, id uuid.UUID, blockchainHash string) error
}

type BlockchainTransactionRepository interface {
	Create(ctx context.Context, transaction *domain.BlockchainTransaction) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.BlockchainTransaction, error)
	GetByTransactionHash(ctx context.Context, hash string) (*domain.BlockchainTransaction, error)
	Update(ctx context.Context, id uuid.UUID, updates map[string]interface{}) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, filter *dto.BlockchainTransactionFilter) ([]*domain.BlockchainTransaction, int64, error)
	GetByEvent(ctx context.Context, eventID uuid.UUID) ([]*domain.BlockchainTransaction, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status string, blockNumber *int64) error
}
