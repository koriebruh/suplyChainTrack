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

type supplyChainService struct {
	repo            repository.SupplyChainEventRepository
	productRepo     repository.ProductRepository
	stakeholderRepo repository.StakeholderRepository
}

func NewSupplyChainService(repo repository.SupplyChainEventRepository, productRepo repository.ProductRepository, stakeholderRepo repository.StakeholderRepository) *supplyChainService {
	return &supplyChainService{repo: repo, productRepo: productRepo, stakeholderRepo: stakeholderRepo}
}

func (s *supplyChainService) CreateEvent(ctx context.Context, req *dto.CreateSupplyChainEventRequest) (*domain.SupplyChainEvent, error) {
	// Validate event type
	if !domain.IsValidEventType(req.EventType) {
		return nil, ErrInvalidEventType
	}

	// Validate product if provided
	if req.ProductID != nil {
		if _, err := s.productRepo.GetByID(ctx, *req.ProductID); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, ErrProductNotFound
			}
			return nil, fmt.Errorf("failed to validate product: %w", err)
		}
	}

	// Validate stakeholder if provided
	if req.StakeholderID != nil {
		if _, err := s.stakeholderRepo.GetByID(ctx, *req.StakeholderID); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, ErrStakeholderNotFound
			}
			return nil, fmt.Errorf("failed to validate stakeholder: %w", err)
		}
	}

	// Validate event sequence
	if err := s.ValidateEventSequence(ctx, req); err != nil {
		return nil, err
	}

	event := &domain.SupplyChainEvent{
		ID:             uuid.New(),
		ProductID:      req.ProductID,
		StakeholderID:  req.StakeholderID,
		EventType:      req.EventType,
		Location:       req.Location,
		Timestamp:      req.Timestamp,
		Metadata:       req.Metadata,
		BlockchainHash: req.BlockchainHash,
		IsVerified:     false,
		CreatedAt:      time.Now(),
	}

	if err := s.repo.Create(ctx, event); err != nil {
		return nil, fmt.Errorf("failed to create supply chain event: %w", err)
	}

	return s.repo.GetByID(ctx, event.ID)
}

func (s *supplyChainService) GetEvent(ctx context.Context, id uuid.UUID) (*domain.SupplyChainEvent, error) {
	event, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrEventNotFound
		}
		return nil, fmt.Errorf("failed to get supply chain event: %w", err)
	}
	return event, nil
}

func (s *supplyChainService) UpdateEvent(ctx context.Context, id uuid.UUID, updates map[string]interface{}) (*domain.SupplyChainEvent, error) {
	if _, err := s.GetEvent(ctx, id); err != nil {
		return nil, err
	}

	if err := s.repo.Update(ctx, id, updates); err != nil {
		return nil, fmt.Errorf("failed to update supply chain event: %w", err)
	}

	return s.GetEvent(ctx, id)
}

func (s *supplyChainService) DeleteEvent(ctx context.Context, id uuid.UUID) error {
	if _, err := s.GetEvent(ctx, id); err != nil {
		return err
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete supply chain event: %w", err)
	}

	return nil
}

func (s *supplyChainService) ListEvents(ctx context.Context, filter *dto.SupplyChainEventFilter) (*dto.PaginatedResponse, error) {
	if filter == nil {
		filter = &dto.SupplyChainEventFilter{Limit: 10, Offset: 0}
	}

	events, total, err := s.repo.List(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to list supply chain events: %w", err)
	}

	return &dto.PaginatedResponse{
		Data:    events,
		Total:   int(total),
		Limit:   filter.Limit,
		Offset:  filter.Offset,
		HasMore: filter.Offset+filter.Limit < int(total),
	}, nil
}

func (s *supplyChainService) GetProductTrace(ctx context.Context, productID uuid.UUID) (*dto.SupplyChainTrace, error) {
	if _, err := s.productRepo.GetByID(ctx, productID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrProductNotFound
		}
		return nil, fmt.Errorf("failed to validate product: %w", err)
	}

	trace, err := s.repo.GetTrace(ctx, productID)
	if err != nil {
		return nil, fmt.Errorf("failed to get product trace: %w", err)
	}

	return trace, nil
}

func (s *supplyChainService) VerifyEvent(ctx context.Context, id uuid.UUID, blockchainHash string) error {
	if _, err := s.GetEvent(ctx, id); err != nil {
		return err
	}

	if err := s.repo.VerifyEvent(ctx, id, blockchainHash); err != nil {
		return fmt.Errorf("failed to verify event: %w", err)
	}

	return nil
}

func (s *supplyChainService) GetEventsByProduct(ctx context.Context, productID uuid.UUID) ([]*domain.SupplyChainEvent, error) {
	if _, err := s.productRepo.GetByID(ctx, productID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrProductNotFound
		}
		return nil, fmt.Errorf("failed to validate product: %w", err)
	}

	events, err := s.repo.GetByProduct(ctx, productID)
	if err != nil {
		return nil, fmt.Errorf("failed to get events by product: %w", err)
	}

	return events, nil
}

func (s *supplyChainService) GetEventsByStakeholder(ctx context.Context, stakeholderID uuid.UUID) ([]*domain.SupplyChainEvent, error) {
	if _, err := s.stakeholderRepo.GetByID(ctx, stakeholderID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrStakeholderNotFound
		}
		return nil, fmt.Errorf("failed to validate stakeholder: %w", err)
	}

	events, err := s.repo.GetByStakeholder(ctx, stakeholderID)
	if err != nil {
		return nil, fmt.Errorf("failed to get events by stakeholder: %w", err)
	}

	return events, nil
}

func (s *supplyChainService) ValidateEventSequence(ctx context.Context, req *dto.CreateSupplyChainEventRequest) error {
	if req.ProductID == nil {
		return nil // Skip validation if no product specified
	}

	// Get existing events for the product
	existingEvents, err := s.repo.GetByProduct(ctx, *req.ProductID)
	if err != nil {
		return fmt.Errorf("failed to get existing events: %w", err)
	}

	// Validate event sequence based on business rules
	switch req.EventType {
	case domain.EventTypeManufactured:
		// Manufactured should be the first event
		for _, event := range existingEvents {
			if event.EventType == domain.EventTypeManufactured {
				return ErrInvalidEventSequence
			}
		}
	case domain.EventTypeShipped:
		// Shipped requires previous manufactured or received event
		hasValidPrevious := false
		for _, event := range existingEvents {
			if event.EventType == domain.EventTypeManufactured || event.EventType == domain.EventTypeReceived {
				hasValidPrevious = true
				break
			}
		}
		if !hasValidPrevious {
			return ErrInvalidEventSequence
		}
	case domain.EventTypeReceived:
		// Received requires previous shipped event
		hasShipped := false
		for _, event := range existingEvents {
			if event.EventType == domain.EventTypeShipped {
				hasShipped = true
				break
			}
		}
		if !hasShipped {
			return ErrInvalidEventSequence
		}
	case domain.EventTypeSold:
		// Sold requires product to be received
		hasReceived := false
		for _, event := range existingEvents {
			if event.EventType == domain.EventTypeReceived {
				hasReceived = true
				break
			}
		}
		if !hasReceived {
			return ErrInvalidEventSequence
		}
	}

	return nil
}
