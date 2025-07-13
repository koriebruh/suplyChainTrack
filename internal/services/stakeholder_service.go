package services

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/koriebruh/suplyChainTrack/internal/domain"
	"github.com/koriebruh/suplyChainTrack/internal/dto"
	"github.com/koriebruh/suplyChainTrack/internal/repository"
	"gorm.io/gorm"
	"time"
)

type stakeholderService struct {
	repo repository.StakeholderRepository
}

func NewStakeholderService(repo repository.StakeholderRepository) *stakeholderService {
	return &stakeholderService{repo: repo}
}

func (s *stakeholderService) CreateStakeholder(ctx context.Context, req *dto.CreateStakeholderRequest) (*domain.Stakeholder, error) {
	if !domain.IsValidStakeholderType(req.Type) {
		return nil, ErrInvalidStakeholderType
	}

	// Check if email already exists
	if _, err := s.repo.GetByEmail(ctx, req.Email); err == nil {
		return nil, ErrDuplicateEmail
	}

	// Check if wallet address already exists (if provided)
	if req.WalletAddress != nil {
		if _, err := s.repo.GetByWalletAddress(ctx, *req.WalletAddress); err == nil {
			return nil, ErrDuplicateWallet
		}
	}

	stakeholder := &domain.Stakeholder{
		ID:            uuid.New(),
		Name:          req.Name,
		Type:          req.Type,
		WalletAddress: req.WalletAddress,
		Email:         req.Email,
		Phone:         req.Phone,
		Address:       req.Address,
		IsVerified:    false,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if err := s.repo.Create(ctx, stakeholder); err != nil {
		return nil, fmt.Errorf("failed to create stakeholder: %w", err)
	}

	return stakeholder, nil
}

func (s *stakeholderService) GetStakeholder(ctx context.Context, id uuid.UUID) (*domain.Stakeholder, error) {
	stakeholder, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrStakeholderNotFound
		}
		return nil, fmt.Errorf("failed to get stakeholder: %w", err)
	}
	return stakeholder, nil
}

func (s *stakeholderService) GetStakeholderByEmail(ctx context.Context, email string) (*domain.Stakeholder, error) {
	stakeholder, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrStakeholderNotFound
		}
		return nil, fmt.Errorf("failed to get stakeholder by email: %w", err)
	}
	return stakeholder, nil
}

func (s *stakeholderService) UpdateStakeholder(ctx context.Context, id uuid.UUID, req *dto.UpdateStakeholderRequest) (*domain.Stakeholder, error) {
	// Check if stakeholder exists
	existing, err := s.GetStakeholder(ctx, id)
	if err != nil {
		return nil, err
	}

	updates := make(map[string]interface{})

	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.Email != nil {
		// Check if new email already exists
		if *req.Email != existing.Email {
			if _, err := s.repo.GetByEmail(ctx, *req.Email); err == nil {
				return nil, ErrDuplicateEmail
			}
		}
		updates["email"] = *req.Email
	}
	if req.WalletAddress != nil {
		// Check if new wallet address already exists
		if existing.WalletAddress == nil || *req.WalletAddress != *existing.WalletAddress {
			if _, err := s.repo.GetByWalletAddress(ctx, *req.WalletAddress); err == nil {
				return nil, ErrDuplicateWallet
			}
		}
		updates["wallet_address"] = *req.WalletAddress
	}
	if req.Phone != nil {
		updates["phone"] = *req.Phone
	}
	if req.Address != nil {
		updates["address"] = *req.Address
	}
	if req.IsVerified != nil {
		updates["is_verified"] = *req.IsVerified
	}

	updates["updated_at"] = time.Now()

	if err := s.repo.Update(ctx, id, updates); err != nil {
		return nil, fmt.Errorf("failed to update stakeholder: %w", err)
	}

	return s.GetStakeholder(ctx, id)
}

func (s *stakeholderService) DeleteStakeholder(ctx context.Context, id uuid.UUID) error {
	if _, err := s.GetStakeholder(ctx, id); err != nil {
		return err
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete stakeholder: %w", err)
	}

	return nil
}

func (s *stakeholderService) ListStakeholders(ctx context.Context, filter *dto.StakeholderFilter) (*dto.PaginatedResponse, error) {
	if filter == nil {
		filter = &dto.StakeholderFilter{Limit: 10, Offset: 0}
	}

	stakeholders, total, err := s.repo.List(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to list stakeholders: %w", err)
	}

	return &dto.PaginatedResponse{
		Data:    stakeholders,
		Total:   int(total),
		Limit:   filter.Limit,
		Offset:  filter.Offset,
		HasMore: filter.Offset+filter.Limit < int(total),
	}, nil
}

func (s *stakeholderService) GetStakeholderStats(ctx context.Context, id uuid.UUID) (*dto.StakeholderStats, error) {
	if _, err := s.GetStakeholder(ctx, id); err != nil {
		return nil, err
	}

	stats, err := s.repo.GetStats(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get stakeholder stats: %w", err)
	}

	return stats, nil
}

func (s *stakeholderService) VerifyStakeholder(ctx context.Context, id uuid.UUID) error {
	updates := map[string]interface{}{
		"is_verified": true,
		"updated_at":  time.Now(),
	}

	if err := s.repo.Update(ctx, id, updates); err != nil {
		return fmt.Errorf("failed to verify stakeholder: %w", err)
	}

	return nil
}
