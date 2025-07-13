package repository

import (
	"context"
	"github.com/google/uuid"
	"github.com/koriebruh/suplyChainTrack/internal/domain"
	"github.com/koriebruh/suplyChainTrack/internal/dto"
	"gorm.io/gorm"
)

type supplyChainEventRepository struct {
	db *gorm.DB
}

func NewSupplyChainEventRepository(db *gorm.DB) *supplyChainEventRepository {
	return &supplyChainEventRepository{db: db}
}

func (r *supplyChainEventRepository) Create(ctx context.Context, event *domain.SupplyChainEvent) error {
	return r.db.WithContext(ctx).Create(event).Error
}

func (r *supplyChainEventRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.SupplyChainEvent, error) {
	var event domain.SupplyChainEvent
	err := r.db.WithContext(ctx).Preload("Product").Preload("Stakeholder").Where("id = ?", id).First(&event).Error
	if err != nil {
		return nil, err
	}
	return &event, nil
}

func (r *supplyChainEventRepository) Update(ctx context.Context, id uuid.UUID, updates map[string]interface{}) error {
	return r.db.WithContext(ctx).Model(&domain.SupplyChainEvent{}).Where("id = ?", id).Updates(updates).Error
}

func (r *supplyChainEventRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&domain.SupplyChainEvent{}, id).Error
}

func (r *supplyChainEventRepository) List(ctx context.Context, filter *dto.SupplyChainEventFilter) ([]*domain.SupplyChainEvent, int64, error) {
	var events []*domain.SupplyChainEvent
	var total int64

	query := r.db.WithContext(ctx).Model(&domain.SupplyChainEvent{}).Preload("Product").Preload("Stakeholder")

	// Apply filters
	if filter.ProductID != nil {
		query = query.Where("product_id = ?", *filter.ProductID)
	}
	if filter.StakeholderID != nil {
		query = query.Where("stakeholder_id = ?", *filter.StakeholderID)
	}
	if filter.EventType != nil {
		query = query.Where("event_type = ?", *filter.EventType)
	}
	if filter.Location != nil {
		query = query.Where("location ILIKE ?", "%"+*filter.Location+"%")
	}
	if filter.IsVerified != nil {
		query = query.Where("is_verified = ?", *filter.IsVerified)
	}
	if filter.FromDate != nil {
		query = query.Where("timestamp >= ?", *filter.FromDate)
	}
	if filter.ToDate != nil {
		query = query.Where("timestamp <= ?", *filter.ToDate)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination and ordering
	query = query.Order("timestamp DESC")
	if filter.Limit > 0 {
		query = query.Limit(filter.Limit)
	}
	if filter.Offset > 0 {
		query = query.Offset(filter.Offset)
	}

	err := query.Find(&events).Error
	return events, total, err
}

func (r *supplyChainEventRepository) GetByProduct(ctx context.Context, productID uuid.UUID) ([]*domain.SupplyChainEvent, error) {
	var events []*domain.SupplyChainEvent
	err := r.db.WithContext(ctx).Preload("Product").Preload("Stakeholder").Where("product_id = ?", productID).Order("timestamp ASC").Find(&events).Error
	return events, err
}

func (r *supplyChainEventRepository) GetByStakeholder(ctx context.Context, stakeholderID uuid.UUID) ([]*domain.SupplyChainEvent, error) {
	var events []*domain.SupplyChainEvent
	err := r.db.WithContext(ctx).Preload("Product").Preload("Stakeholder").Where("stakeholder_id = ?", stakeholderID).Order("timestamp DESC").Find(&events).Error
	return events, err
}

func (r *supplyChainEventRepository) GetTrace(ctx context.Context, productID uuid.UUID) (*dto.SupplyChainTrace, error) {
	// Get product
	var product domain.Product
	if err := r.db.WithContext(ctx).Preload("Manufacturer").Where("id = ?", productID).First(&product).Error; err != nil {
		return nil, err
	}

	// Get all events for the product
	events, err := r.GetByProduct(ctx, productID)
	if err != nil {
		return nil, err
	}

	return &dto.SupplyChainTrace{
		Product: &product,
		Events:  events,
	}, nil
}

func (r *supplyChainEventRepository) VerifyEvent(ctx context.Context, id uuid.UUID, blockchainHash string) error {
	updates := map[string]interface{}{
		"is_verified":     true,
		"blockchain_hash": blockchainHash,
	}
	return r.db.WithContext(ctx).Model(&domain.SupplyChainEvent{}).Where("id = ?", id).Updates(updates).Error
}
