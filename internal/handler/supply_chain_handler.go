package handler

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/koriebruh/suplyChainTrack/internal/dto"
	"github.com/koriebruh/suplyChainTrack/internal/services"
	"strconv"
)

type stakeHolderHandler struct {
	service services.StakeholderService
}

func (h *stakeHolderHandler) CreateStakeholder(c *fiber.Ctx) error {
	var req dto.CreateStakeholderRequest
	if err := c.BodyParser(&req); err != nil {
		return SendError(c, fiber.StatusBadRequest, err, "Invalid request body")
	}

	stakeholder, err := h.service.CreateStakeholder(c.Context(), &req)
	if err != nil {
		switch err {
		case services.ErrDuplicateEmail:
			return SendError(c, fiber.StatusConflict, err, "Email already exists")
		case services.ErrDuplicateWallet:
			return SendError(c, fiber.StatusConflict, err, "Wallet address already exists")
		case services.ErrInvalidStakeholderType:
			return SendError(c, fiber.StatusBadRequest, err, "Invalid stakeholder type")
		default:
			return SendError(c, fiber.StatusInternalServerError, err, "Failed to create stakeholder")
		}
	}

	return SendSuccess(c, fiber.StatusCreated, stakeholder, "Stakeholder created successfully")
}

func (h *stakeHolderHandler) GetStakeholderByEmail(c *fiber.Ctx) error {
	email := c.Query("email")
	if email == "" {
		return SendError(c, fiber.StatusBadRequest, fmt.Errorf("email parameter required"), "Email parameter is required")
	}

	stakeholder, err := h.service.GetStakeholderByEmail(c.Context(), email)
	if err != nil {
		if err == services.ErrStakeholderNotFound {
			return SendError(c, fiber.StatusNotFound, err, "Stakeholder not found")
		}
		return SendError(c, fiber.StatusInternalServerError, err, "Failed to get stakeholder")
	}

	return SendSuccess(c, fiber.StatusOK, stakeholder, "Stakeholder retrieved successfully")
}

func (h *stakeHolderHandler) GetStakeholder(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return SendError(c, fiber.StatusBadRequest, err, "Invalid stakeholder ID")
	}

	stakeholder, err := h.service.GetStakeholder(c.Context(), id)
	if err != nil {
		if err == services.ErrStakeholderNotFound {
			return SendError(c, fiber.StatusNotFound, err, "Stakeholder not found")
		}
		return SendError(c, fiber.StatusInternalServerError, err, "Failed to get stakeholder")
	}

	return SendSuccess(c, fiber.StatusOK, stakeholder, "Stakeholder retrieved successfully")
}

func (h *stakeHolderHandler) UpdateStakeholder(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return SendError(c, fiber.StatusBadRequest, err, "Invalid stakeholder ID")
	}

	var req dto.UpdateStakeholderRequest
	if err := c.BodyParser(&req); err != nil {
		return SendError(c, fiber.StatusBadRequest, err, "Invalid request body")
	}

	stakeholder, err := h.service.UpdateStakeholder(c.Context(), id, &req)
	if err != nil {
		switch err {
		case services.ErrStakeholderNotFound:
			return SendError(c, fiber.StatusNotFound, err, "Stakeholder not found")
		case services.ErrDuplicateEmail:
			return SendError(c, fiber.StatusConflict, err, "Email already exists")
		case services.ErrDuplicateWallet:
			return SendError(c, fiber.StatusConflict, err, "Wallet address already exists")
		default:
			return SendError(c, fiber.StatusInternalServerError, err, "Failed to update stakeholder")
		}
	}

	return SendSuccess(c, fiber.StatusOK, stakeholder, "Stakeholder updated successfully")
}

func (h *stakeHolderHandler) DeleteStakeholder(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return SendError(c, fiber.StatusBadRequest, err, "Invalid stakeholder ID")
	}

	err = h.service.DeleteStakeholder(c.Context(), id)
	if err != nil {
		if err == services.ErrStakeholderNotFound {
			return SendError(c, fiber.StatusNotFound, err, "Stakeholder not found")
		}
		return SendError(c, fiber.StatusInternalServerError, err, "Failed to delete stakeholder")
	}

	return SendSuccess(c, fiber.StatusOK, nil, "Stakeholder deleted successfully")
}

func (h *stakeHolderHandler) ListStakeholders(c *fiber.Ctx) error {
	filter := &dto.StakeholderFilter{}

	// Parse query parameters
	if limit := c.Query("limit"); limit != "" {
		if l, err := strconv.Atoi(limit); err == nil {
			filter.Limit = l
		}
	}
	if offset := c.Query("offset"); offset != "" {
		if o, err := strconv.Atoi(offset); err == nil {
			filter.Offset = o
		}
	}
	if stakeholderType := c.Query("type"); stakeholderType != "" {
		filter.Type = &stakeholderType
	}
	if isVerified := c.Query("is_verified"); isVerified != "" {
		if v, err := strconv.ParseBool(isVerified); err == nil {
			filter.IsVerified = &v
		}
	}
	if email := c.Query("email"); email != "" {
		filter.Email = &email
	}

	// Set default values
	if filter.Limit == 0 {
		filter.Limit = 10
	}

	response, err := h.service.ListStakeholders(c.Context(), filter)
	if err != nil {
		return SendError(c, fiber.StatusInternalServerError, err, "Failed to list stakeholders")
	}

	return SendSuccess(c, fiber.StatusOK, response, "Stakeholders retrieved successfully")
}

func (h *stakeHolderHandler) GetStakeholderStats(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return SendError(c, fiber.StatusBadRequest, err, "Invalid stakeholder ID")
	}

	stats, err := h.service.GetStakeholderStats(c.Context(), id)
	if err != nil {
		if err == services.ErrStakeholderNotFound {
			return SendError(c, fiber.StatusNotFound, err, "Stakeholder not found")
		}
		return SendError(c, fiber.StatusInternalServerError, err, "Failed to get stakeholder stats")
	}

	return SendSuccess(c, fiber.StatusOK, stats, "Stakeholder stats retrieved successfully")
}

func (h *stakeHolderHandler) VerifyStakeholder(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return SendError(c, fiber.StatusBadRequest, err, "Invalid stakeholder ID")
	}

	err = h.service.VerifyStakeholder(c.Context(), id)
	if err != nil {
		if err == services.ErrStakeholderNotFound {
			return SendError(c, fiber.StatusNotFound, err, "Stakeholder not found")
		}
		return SendError(c, fiber.StatusInternalServerError, err, "Failed to verify stakeholder")
	}

	return SendSuccess(c, fiber.StatusOK, nil, "Stakeholder verified successfully")
}
