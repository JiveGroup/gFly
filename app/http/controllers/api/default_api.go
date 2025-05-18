package api

import (
	"fmt"
	"gfly/app/http/response"
	"github.com/gflydev/core"
	"github.com/gflydev/core/utils"
)

// ====================================================================
// ======================== Controller Creation =======================
// ====================================================================

// NewDefaultApi As a constructor to create new API.
func NewDefaultApi() *DefaultApi {
	return &DefaultApi{}
}

// DefaultApi API struct.
type DefaultApi struct {
	core.Api
}

// ====================================================================
// ========================= Request Handling =========================
// ====================================================================

// Handle Process main logic for API.
// @Summary Get API info
// @Description Get API server information
// @Tags Misc
// @Accept json
// @Produce json
// @Success 200 {object} response.ServerInfo
// @Router /info [get]
func (h *DefaultApi) Handle(c *core.Ctx) error {
	obj := response.ServerInfo{
		Name: utils.Getenv("API_NAME", "gfly"),
		Prefix: fmt.Sprintf(
			"/%s/%s",
			utils.Getenv("API_PREFIX", "api"),
			utils.Getenv("API_VERSION", "v1"),
		),
		Server: utils.Getenv("APP_URL", "http://localhost:7789"),
	}

	return c.JSONAny(obj)
}
