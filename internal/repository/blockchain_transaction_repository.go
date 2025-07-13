package repository

import (
	"context"
	"github.com/google/uuid"
	"github.com/koriebruh/suplyChainTrack/internal/domain"
	"github.com/koriebruh/suplyChainTrack/internal/dto"
	"gorm.io/gorm"
)

type blockchainTransactionRepository struct {
	db *gorm.DB
}

func NewBlockchainTransactionRepository(db *gorm.DB) *blockchainTransactionRepository {
	return &blockchainTransactionRepository{db: db}
}

func (r *blockchainTransactionRepository) Create(ctx context.Context, transaction *domain.BlockchainTransaction) error {
	return r.db.WithContext(ctx).Create(transaction).Error
}

func (r *blockchainTransactionRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.BlockchainTransaction, error) {
	var transaction domain.BlockchainTransaction
	err := r.db.WithContext(ctx).Preload("Event").Where("id = ?", id).First(&transaction).Error
	if err != nil {
		return nil, err
	}
	return &transaction, nil
}

func (r *blockchainTransactionRepository) GetByTransactionHash(ctx context.Context, hash string) (*domain.BlockchainTransaction, error) {
	var transaction domain.BlockchainTransaction
	err := r.db.WithContext(ctx).Preload("Event").Where("transaction_hash = ?", hash).First(&transaction).Error
	if err != nil {
		return nil, err
	}
	return &transaction, nil
}

func (r *blockchainTransactionRepository) Update(ctx context.Context, id uuid.UUID, updates map[string]interface{}) error {
	return r.db.WithContext(ctx).Model(&domain.BlockchainTransaction{}).Where("id = ?", id).Updates(updates).Error
}

func (r *blockchainTransactionRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&domain.BlockchainTransaction{}, id).Error
}

func (r *blockchainTransactionRepository) List(ctx context.Context, filter *dto.BlockchainTransactionFilter) ([]*domain.BlockchainTransaction, int64, error) {
	var transactions []*domain.BlockchainTransaction
	var total int64

	query := r.db.WithContext(ctx).Model(&domain.BlockchainTransaction{}).Preload("Event")

	// Apply filters
	if filter.EventID != nil {
		query = query.Where("event_id = ?", *filter.EventID)
	}
	if filter.Status != nil {
		query = query.Where("status = ?", *filter.Status)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination and ordering
	query = query.Order("created_at DESC")
	if filter.Limit > 0 {
		query = query.Limit(filter.Limit)
	}
	if filter.Offset > 0 {
		query = query.Offset(filter.Offset)
	}

	err := query.Find(&transactions).Error
	return transactions, total, err
}

func (r *blockchainTransactionRepository) GetByEvent(ctx context.Context, eventID uuid.UUID) ([]*domain.BlockchainTransaction, error) {
	var transactions []*domain.BlockchainTransaction
	err := r.db.WithContext(ctx).Preload("Event").Where("event_id = ?", eventID).Order("created_at DESC").Find(&transactions).Error
	return transactions, err
}

func (r *blockchainTransactionRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status string, blockNumber *int64) error {
	updates := map[string]interface{}{
		"status": status,
	}
	if blockNumber != nil {
		updates["block_number"] = *blockNumber
	}
	return r.db.WithContext(ctx).Model(&domain.BlockchainTransaction{}).Where("id = ?", id).Updates(updates).Error
}
