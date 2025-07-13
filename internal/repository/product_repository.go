package repository

import (
	"context"
	"github.com/google/uuid"
	"github.com/koriebruh/suplyChainTrack/internal/domain"
	"github.com/koriebruh/suplyChainTrack/internal/dto"
	"gorm.io/gorm"
)

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) *productRepository {
	return &productRepository{db: db}
}

func (r *productRepository) Create(ctx context.Context, product *domain.Product) error {
	return r.db.WithContext(ctx).Create(product).Error
}

func (r *productRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Product, error) {
	var product domain.Product
	err := r.db.WithContext(ctx).Preload("Manufacturer").Where("id = ?", id).First(&product).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *productRepository) GetBySKU(ctx context.Context, sku string) (*domain.Product, error) {
	var product domain.Product
	err := r.db.WithContext(ctx).Preload("Manufacturer").Where("sku = ?", sku).First(&product).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *productRepository) Update(ctx context.Context, id uuid.UUID, updates map[string]interface{}) error {
	return r.db.WithContext(ctx).Model(&domain.Product{}).Where("id = ?", id).Updates(updates).Error
}

func (r *productRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&domain.Product{}, id).Error
}

func (r *productRepository) List(ctx context.Context, filter *dto.ProductFilter) ([]*domain.Product, int64, error) {
	var products []*domain.Product
	var total int64

	query := r.db.WithContext(ctx).Model(&domain.Product{}).Preload("Manufacturer")

	// Apply filters
	if filter.Category != nil {
		query = query.Where("category = ?", *filter.Category)
	}
	if filter.ManufacturerID != nil {
		query = query.Where("manufacturer_id = ?", *filter.ManufacturerID)
	}
	if filter.SKU != nil {
		query = query.Where("sku ILIKE ?", "%"+*filter.SKU+"%")
	}
	if filter.Name != nil {
		query = query.Where("name ILIKE ?", "%"+*filter.Name+"%")
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	if filter.Limit > 0 {
		query = query.Limit(filter.Limit)
	}
	if filter.Offset > 0 {
		query = query.Offset(filter.Offset)
	}

	err := query.Find(&products).Error
	return products, total, err
}

func (r *productRepository) GetStats(ctx context.Context, id uuid.UUID) (*dto.ProductStats, error) {
	var stats dto.ProductStats
	stats.ProductID = id

	// Get total events
	r.db.WithContext(ctx).Model(&domain.SupplyChainEvent{}).Where("product_id = ?", id).Count(&stats.TotalEvents)

	// Get verified events
	r.db.WithContext(ctx).Model(&domain.SupplyChainEvent{}).Where("product_id = ? AND is_verified = ?", id, true).Count(&stats.VerifiedEvents)

	// Get last event for location and stakeholder
	var lastEvent domain.SupplyChainEvent
	if err := r.db.WithContext(ctx).Where("product_id = ?", id).Order("timestamp DESC").First(&lastEvent).Error; err == nil {
		stats.CurrentLocation = lastEvent.Location
		stats.LastStakeholder = lastEvent.StakeholderID
		stats.LastActivity = lastEvent.Timestamp
	}

	return &stats, nil
}

func (r *productRepository) GetByManufacturer(ctx context.Context, manufacturerID uuid.UUID) ([]*domain.Product, error) {
	var products []*domain.Product
	err := r.db.WithContext(ctx).Preload("Manufacturer").Where("manufacturer_id = ?", manufacturerID).Find(&products).Error
	return products, err
}
