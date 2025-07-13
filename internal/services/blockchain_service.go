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

type blockchainService struct {
	repo      repository.BlockchainTransactionRepository
	eventRepo repository.SupplyChainEventRepository
}

func NewBlockchainService(repo repository.BlockchainTransactionRepository, eventRepo repository.SupplyChainEventRepository) *blockchainService {
	return &blockchainService{repo: repo, eventRepo: eventRepo}
}

func (s *blockchainService) CreateTransaction(ctx context.Context, req *dto.CreateBlockchainTransactionRequest) (*domain.BlockchainTransaction, error) {
	// Validate transaction status
	if !domain.IsValidTransactionStatus(req.Status) {
		return nil, ErrInvalidTransactionStatus
	}

	// Validate event if provided
	if req.EventID != nil {
		if _, err := s.eventRepo.GetByID(ctx, *req.EventID); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, ErrEventNotFound
			}
			return nil, fmt.Errorf("failed to validate event: %w", err)
		}
	}

	transaction := &domain.BlockchainTransaction{
		ID:              uuid.New(),
		EventID:         req.EventID,
		TransactionHash: req.TransactionHash,
		BlockNumber:     req.BlockNumber,
		GasUsed:         req.GasUsed,
		Status:          req.Status,
		CreatedAt:       time.Now(),
	}

	if err := s.repo.Create(ctx, transaction); err != nil {
		return nil, fmt.Errorf("failed to create blockchain transaction: %w", err)
	}

	return s.repo.GetByID(ctx, transaction.ID)
}

func (s *blockchainService) GetTransaction(ctx context.Context, id uuid.UUID) (*domain.BlockchainTransaction, error) {
	transaction, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTransactionNotFound
		}
		return nil, fmt.Errorf("failed to get blockchain transaction: %w", err)
	}
	return transaction, nil
}

func (s *blockchainService) GetTransactionByHash(ctx context.Context, hash string) (*domain.BlockchainTransaction, error) {
	transaction, err := s.repo.GetByTransactionHash(ctx, hash)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTransactionNotFound
		}
		return nil, fmt.Errorf("failed to get blockchain transaction by hash: %w", err)
	}
	return transaction, nil
}

func (s *blockchainService) UpdateTransaction(ctx context.Context, id uuid.UUID, updates map[string]interface{}) (*domain.BlockchainTransaction, error) {
	if _, err := s.GetTransaction(ctx, id); err != nil {
		return nil, err
	}

	if err := s.repo.Update(ctx, id, updates); err != nil {
		return nil, fmt.Errorf("failed to update blockchain transaction: %w", err)
	}

	return s.GetTransaction(ctx, id)
}

func (s *blockchainService) ListTransactions(ctx context.Context, filter *dto.BlockchainTransactionFilter) (*dto.PaginatedResponse, error) {
	if filter == nil {
		filter = &dto.BlockchainTransactionFilter{Limit: 10, Offset: 0}
	}

	transactions, total, err := s.repo.List(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to list blockchain transactions: %w", err)
	}

	return &dto.PaginatedResponse{
		Data:    transactions,
		Total:   int(total),
		Limit:   filter.Limit,
		Offset:  filter.Offset,
		HasMore: filter.Offset+filter.Limit < int(total),
	}, nil
}

func (s *blockchainService) UpdateTransactionStatus(ctx context.Context, id uuid.UUID, status string, blockNumber *int64) error {
	if !domain.IsValidTransactionStatus(status) {
		return ErrInvalidTransactionStatus
	}

	updates := map[string]interface{}{
		"status":     status,
		"updated_at": time.Now(),
	}

	if blockNumber != nil {
		updates["block_number"] = *blockNumber
	}

	if err := s.repo.Update(ctx, id, updates); err != nil {
		return fmt.Errorf("failed to update transaction status: %w", err)
	}

	return nil
}

func (s *blockchainService) GetTransactionsByEvent(ctx context.Context, eventID uuid.UUID) ([]*domain.BlockchainTransaction, error) {
	// Validate event exists
	if _, err := s.eventRepo.GetByID(ctx, eventID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrEventNotFound
		}
		return nil, fmt.Errorf("failed to validate event: %w", err)
	}

	transactions, err := s.repo.GetByEvent(ctx, eventID)
	if err != nil {
		return nil, fmt.Errorf("failed to get transactions by event: %w", err)
	}

	return transactions, nil
}
