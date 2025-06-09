package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/koriebruh/suplyChainTrack/conf"
	"github.com/koriebruh/suplyChainTrack/internal/dto"
)

type MetricHandler interface {
	Health(c *fiber.Ctx) error
}

type MetricHandlerImpl struct {
	conf.Config
}

func NewMetricHandlerImpl(config conf.Config) *MetricHandlerImpl {
	return &MetricHandlerImpl{Config: config}
}

// Health godoc
// @Summary Show the health status of the service
// @Description Returns service health status and version
// @Tags Health
// @Accept  json
// @Produce  json
// @Success 200 {object} dto.HealthResponse
// @Router /health [get]
func (m MetricHandlerImpl) Health(c *fiber.Ctx) error {
	return c.JSON(dto.HealthResponse{
		Success:   true,
		Status:    "healthy",
		Timestamp: c.Context().Time().UnixNano(),
		Version:   m.Config.AppConfig.Version,
	})
}
