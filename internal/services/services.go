package services

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/koriebruh/suplyChainTrack/internal/domain"
	"github.com/koriebruh/suplyChainTrack/internal/dto"
	"github.com/koriebruh/suplyChainTrack/internal/repository"
)

// custom error definitions for the supply chain tracking service
var (
	ErrStakeholderNotFound      = errors.New("stakeholder not found")
	ErrProductNotFound          = errors.New("product not found")
	ErrEventNotFound            = errors.New("event not found")
	ErrTransactionNotFound      = errors.New("transaction not found")
	ErrDuplicateEmail           = errors.New("email already exists")
	ErrDuplicateSKU             = errors.New("SKU already exists")
	ErrDuplicateWallet          = errors.New("wallet address already exists")
	ErrInvalidStakeholderType   = errors.New("invalid stakeholder type")
	ErrInvalidEventType         = errors.New("invalid event type")
	ErrInvalidTransactionStatus = errors.New("invalid transaction status")
	ErrUnauthorized             = errors.New("unauthorized access")
	ErrInvalidEventSequence     = errors.New("invalid event sequence")
)

type ServiceManager struct {
	Stakeholder StakeholderService
	Product     ProductService
	SupplyChain SupplyChainService
	Blockchain  BlockchainService
}

func NewServiceManager(repos *repository.RepositoriesManagers) *ServiceManager {
	return &ServiceManager{
		Stakeholder: NewStakeholderService(repos.Stakeholder),
		Product:     NewProductService(repos.Product, repos.Stakeholder),
		SupplyChain: NewSupplyChainService(repos.SupplyChainEvent, repos.Product, repos.Stakeholder),
		Blockchain:  NewBlockchainService(repos.BlockchainTransaction, repos.SupplyChainEvent),
	}
}

type StakeholderService interface {
	CreateStakeholder(ctx context.Context, req *dto.CreateStakeholderRequest) (*domain.Stakeholder, error)
	GetStakeholder(ctx context.Context, id uuid.UUID) (*domain.Stakeholder, error)
	GetStakeholderByEmail(ctx context.Context, email string) (*domain.Stakeholder, error)
	UpdateStakeholder(ctx context.Context, id uuid.UUID, req *dto.UpdateStakeholderRequest) (*domain.Stakeholder, error)
	DeleteStakeholder(ctx context.Context, id uuid.UUID) error
	ListStakeholders(ctx context.Context, filter *dto.StakeholderFilter) (*dto.PaginatedResponse, error)
	GetStakeholderStats(ctx context.Context, id uuid.UUID) (*dto.StakeholderStats, error)
	VerifyStakeholder(ctx context.Context, id uuid.UUID) error
}

type ProductService interface {
	CreateProduct(ctx context.Context, req *dto.CreateProductRequest) (*domain.Product, error)
	GetProduct(ctx context.Context, id uuid.UUID) (*domain.Product, error)
	GetProductBySKU(ctx context.Context, sku string) (*domain.Product, error)
	UpdateProduct(ctx context.Context, id uuid.UUID, req *dto.UpdateProductRequest) (*domain.Product, error)
	DeleteProduct(ctx context.Context, id uuid.UUID) error
	ListProducts(ctx context.Context, filter *dto.ProductFilter) (*dto.PaginatedResponse, error)
	GetProductStats(ctx context.Context, id uuid.UUID) (*dto.ProductStats, error)
	GetProductsByManufacturer(ctx context.Context, manufacturerID uuid.UUID) ([]*domain.Product, error)
}

type SupplyChainService interface {
	CreateEvent(ctx context.Context, req *dto.CreateSupplyChainEventRequest) (*domain.SupplyChainEvent, error)
	GetEvent(ctx context.Context, id uuid.UUID) (*domain.SupplyChainEvent, error)
	UpdateEvent(ctx context.Context, id uuid.UUID, updates map[string]interface{}) (*domain.SupplyChainEvent, error)
	DeleteEvent(ctx context.Context, id uuid.UUID) error
	ListEvents(ctx context.Context, filter *dto.SupplyChainEventFilter) (*dto.PaginatedResponse, error)
	GetProductTrace(ctx context.Context, productID uuid.UUID) (*dto.SupplyChainTrace, error)
	VerifyEvent(ctx context.Context, id uuid.UUID, blockchainHash string) error
	GetEventsByProduct(ctx context.Context, productID uuid.UUID) ([]*domain.SupplyChainEvent, error)
	GetEventsByStakeholder(ctx context.Context, stakeholderID uuid.UUID) ([]*domain.SupplyChainEvent, error)
	ValidateEventSequence(ctx context.Context, req *dto.CreateSupplyChainEventRequest) error
}

type BlockchainService interface {
	CreateTransaction(ctx context.Context, req *dto.CreateBlockchainTransactionRequest) (*domain.BlockchainTransaction, error)
	GetTransaction(ctx context.Context, id uuid.UUID) (*domain.BlockchainTransaction, error)
	GetTransactionByHash(ctx context.Context, hash string) (*domain.BlockchainTransaction, error)
	UpdateTransaction(ctx context.Context, id uuid.UUID, updates map[string]interface{}) (*domain.BlockchainTransaction, error)
	ListTransactions(ctx context.Context, filter *dto.BlockchainTransactionFilter) (*dto.PaginatedResponse, error)
	UpdateTransactionStatus(ctx context.Context, id uuid.UUID, status string, blockNumber *int64) error
	GetTransactionsByEvent(ctx context.Context, eventID uuid.UUID) ([]*domain.BlockchainTransaction, error)
}
