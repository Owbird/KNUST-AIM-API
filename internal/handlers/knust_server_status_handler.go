package handlers

import (
	"net/http"

	"github.com/Owbird/KNUST-AIM-API/models"
	"github.com/Owbird/KNUST-AIM-API/pkg/status"
	"github.com/gin-gonic/gin"
)

var statusFunctions = status.NewStatusFunctions()

// @Summary Get the status of KNUST servers
// @Description This checks which of the used KNUST servers are up or down
// @Tags KNUST Servers
// @Produce json
// @Success 200 {object} models.KNUSTServerStatusResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /knust-server-status [get]
func (h *Handlers) KNUSTServerStatusHandler(c *gin.Context) {
	servers := statusFunctions.GetKNUSTStatus()

	c.JSON(http.StatusOK, models.KNUSTServerStatusResponse{
		Message: "Fetched server status successfully",
		Servers: servers,
	})
}

// @Summary Get the status of KNUST servers as a badge
// @Description This sums up the status of the servers and returns a url to an SVG badge from shields.io
// @Tags KNUST Servers
// @Success 200 {string} string "OK"
// @Failure 500 {object} models.ErrorResponse
// @Router /knust-server-status/badge [get]
func (h *Handlers) KNUSTServerStatusBadgeHandler(c *gin.Context) {
	badge, err := statusFunctions.GetStatusBadge()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Message: "Couldn't fetch badge"})
	}

	c.String(http.StatusOK, badge)
}
