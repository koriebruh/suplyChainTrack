package repository

import (
	"context"
	"github.com/google/uuid"
	"github.com/koriebruh/suplyChainTrack/internal/domain"
	"github.com/koriebruh/suplyChainTrack/internal/dto"
	"gorm.io/gorm"
)

type stakeholderRepository struct {
	db *gorm.DB
}

func NewStakeholderRepository(db *gorm.DB) *stakeholderRepository {
	return &stakeholderRepository{db: db}
}

func (r stakeholderRepository) Create(ctx context.Context, stakeholder *domain.Stakeholder) error {
	return r.db.WithContext(ctx).Create(stakeholder).Error
}

func (r stakeholderRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Stakeholder, error) {
	var stakeholder domain.Stakeholder
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&stakeholder).Error
	if err != nil {
		return nil, err
	}
	return &stakeholder, nil
}

func (r stakeholderRepository) GetByEmail(ctx context.Context, email string) (*domain.Stakeholder, error) {
	var stakeholder domain.Stakeholder
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&stakeholder).Error
	if err != nil {
		return nil, err
	}
	return &stakeholder, nil
}

func (r stakeholderRepository) GetByWalletAddress(ctx context.Context, address string) (*domain.Stakeholder, error) {
	var stakeholder domain.Stakeholder
	err := r.db.WithContext(ctx).Where("wallet_address = ?", address).First(&stakeholder).Error
	if err != nil {
		return nil, err
	}
	return &stakeholder, nil
}

func (r stakeholderRepository) Update(ctx context.Context, id uuid.UUID, updates map[string]interface{}) error {
	return r.db.WithContext(ctx).Model(&domain.Stakeholder{}).Where("id = ?", id).Updates(updates).Error
}

func (r stakeholderRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&domain.Stakeholder{}, id).Error
}

func (r stakeholderRepository) List(ctx context.Context, filter *dto.StakeholderFilter) ([]*domain.Stakeholder, int64, error) {
	var stakeholders []*domain.Stakeholder
	var total int64

	query := r.db.WithContext(ctx).Model(&domain.Stakeholder{})

	// Apply filters
	if filter.Type != nil {
		query = query.Where("type = ?", *filter.Type)
	}
	if filter.IsVerified != nil {
		query = query.Where("is_verified = ?", *filter.IsVerified)
	}
	if filter.Email != nil {
		query = query.Where("email ILIKE ?", "%"+*filter.Email+"%")
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

	err := query.Find(&stakeholders).Error
	return stakeholders, total, err
}

func (r stakeholderRepository) GetStats(ctx context.Context, id uuid.UUID) (*dto.StakeholderStats, error) {
	var stats dto.StakeholderStats
	stats.StakeholderID = id

	// Get total products (for manufacturers)
	r.db.WithContext(ctx).Model(&domain.Product{}).Where("manufacturer_id = ?", id).Count(&stats.TotalProducts)

	// Get total events
	r.db.WithContext(ctx).Model(&domain.SupplyChainEvent{}).Where("stakeholder_id = ?", id).Count(&stats.TotalEvents)

	// Get verified events
	r.db.WithContext(ctx).Model(&domain.SupplyChainEvent{}).Where("stakeholder_id = ? AND is_verified = ?", id, true).Count(&stats.VerifiedEvents)

	// Get pending events
	r.db.WithContext(ctx).Model(&domain.SupplyChainEvent{}).Where("stakeholder_id = ? AND is_verified = ?", id, false).Count(&stats.PendingEvents)

	// Get last activity
	var lastEvent domain.SupplyChainEvent
	if err := r.db.WithContext(ctx).Where("stakeholder_id = ?", id).Order("created_at DESC").First(&lastEvent).Error; err == nil {
		stats.LastActivity = lastEvent.CreatedAt
	}

	return &stats, nil
}
