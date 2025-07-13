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

type productService struct {
	repo            repository.ProductRepository
	stakeholderRepo repository.StakeholderRepository
}

func NewProductService(repo repository.ProductRepository, stakeholderRepo repository.StakeholderRepository) *productService {
	return &productService{repo: repo, stakeholderRepo: stakeholderRepo}
}

func (s *productService) CreateProduct(ctx context.Context, req *dto.CreateProductRequest) (*domain.Product, error) {
	// Check if SKU already exists
	if _, err := s.repo.GetBySKU(ctx, req.SKU); err == nil {
		return nil, ErrDuplicateSKU
	}

	// Validate manufacturer if provided
	if req.ManufacturerID != nil {
		manufacturer, err := s.stakeholderRepo.GetByID(ctx, *req.ManufacturerID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, ErrStakeholderNotFound
			}
			return nil, fmt.Errorf("failed to validate manufacturer: %w", err)
		}
		if manufacturer.Type != domain.StakeholderTypeManufacturer {
			return nil, ErrInvalidStakeholderType
		}
	}

	product := &domain.Product{
		ID:             uuid.New(),
		SKU:            req.SKU,
		Name:           req.Name,
		Description:    req.Description,
		Category:       req.Category,
		ManufacturerID: req.ManufacturerID,
		Metadata:       req.Metadata,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	if err := s.repo.Create(ctx, product); err != nil {
		return nil, fmt.Errorf("failed to create product: %w", err)
	}

	return s.repo.GetByID(ctx, product.ID)
}

func (s *productService) GetProduct(ctx context.Context, id uuid.UUID) (*domain.Product, error) {
	product, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrProductNotFound
		}
		return nil, fmt.Errorf("failed to get product: %w", err)
	}
	return product, nil
}

func (s *productService) GetProductBySKU(ctx context.Context, sku string) (*domain.Product, error) {
	product, err := s.repo.GetBySKU(ctx, sku)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrProductNotFound
		}
		return nil, fmt.Errorf("failed to get product by SKU: %w", err)
	}
	return product, nil
}

func (s *productService) UpdateProduct(ctx context.Context, id uuid.UUID, req *dto.UpdateProductRequest) (*domain.Product, error) {
	// Check if product exists
	if _, err := s.GetProduct(ctx, id); err != nil {
		return nil, err
	}

	updates := make(map[string]interface{})

	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.Description != nil {
		updates["description"] = *req.Description
	}
	if req.Category != nil {
		updates["category"] = *req.Category
	}
	if req.ManufacturerID != nil {
		// Validate manufacturer
		manufacturer, err := s.stakeholderRepo.GetByID(ctx, *req.ManufacturerID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, ErrStakeholderNotFound
			}
			return nil, fmt.Errorf("failed to validate manufacturer: %w", err)
		}
		if manufacturer.Type != domain.StakeholderTypeManufacturer {
			return nil, ErrInvalidStakeholderType
		}
		updates["manufacturer_id"] = *req.ManufacturerID
	}
	if req.Metadata != nil {
		updates["metadata"] = req.Metadata
	}

	updates["updated_at"] = time.Now()

	if err := s.repo.Update(ctx, id, updates); err != nil {
		return nil, fmt.Errorf("failed to update product: %w", err)
	}

	return s.GetProduct(ctx, id)
}

func (s *productService) DeleteProduct(ctx context.Context, id uuid.UUID) error {
	if _, err := s.GetProduct(ctx, id); err != nil {
		return err
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete product: %w", err)
	}

	return nil
}

func (s *productService) ListProducts(ctx context.Context, filter *dto.ProductFilter) (*dto.PaginatedResponse, error) {
	if filter == nil {
		filter = &dto.ProductFilter{Limit: 10, Offset: 0}
	}

	products, total, err := s.repo.List(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to list products: %w", err)
	}

	return &dto.PaginatedResponse{
		Data:    products,
		Total:   int(total),
		Limit:   filter.Limit,
		Offset:  filter.Offset,
		HasMore: filter.Offset+filter.Limit < int(total),
	}, nil
}

func (s *productService) GetProductStats(ctx context.Context, id uuid.UUID) (*dto.ProductStats, error) {
	if _, err := s.GetProduct(ctx, id); err != nil {
		return nil, err
	}

	stats, err := s.repo.GetStats(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get product stats: %w", err)
	}

	return stats, nil
}

func (s *productService) GetProductsByManufacturer(ctx context.Context, manufacturerID uuid.UUID) ([]*domain.Product, error) {
	manufacturer, err := s.stakeholderRepo.GetByID(ctx, manufacturerID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrStakeholderNotFound
		}
		return nil, fmt.Errorf("failed to validate manufacturer: %w", err)
	}
	if manufacturer.Type != domain.StakeholderTypeManufacturer {
		return nil, ErrInvalidStakeholderType
	}

	products, err := s.repo.GetByManufacturer(ctx, manufacturerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get products by manufacturer: %w", err)
	}

	return products, nil
}
