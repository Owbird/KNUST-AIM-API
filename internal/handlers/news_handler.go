package handlers

import (
	"net/http"

	"github.com/Owbird/KNUST-AIM-API/models"
	"github.com/Owbird/KNUST-AIM-API/pkg/news"
	"github.com/gin-gonic/gin"
)

var newsFunctions = news.NewNewsFunctions()

// @Summary Get latest news
// @Description Returns the latest news available
// @Tags News
// @Produce json
// @Success 200 {object} models.NewsResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /news [get]
func (h *Handlers) GetNewsHandler(c *gin.Context) {
	appNews, err := newsFunctions.GetNews()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Message: "Couldn't fetch news"})
		return
	}

	c.JSON(http.StatusOK, models.NewsResponse{
		Message: "Fetched news successfully",
		News:    appNews,
	})
}

// @Summary Get news post details
// @Description Get the post details of the news based on the slug
// @Tags News
// @Produce json
// @Accept  json
// @Param slug path string true "News slug"
// @Success 200 {object} models.NewsDetailsResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /news/{slug} [get]
func (h *Handlers) GetNewsDetailsHandler(c *gin.Context) {
	slug := c.Param("slug")

	details, err := newsFunctions.GetNewsDetails(slug)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Message: "Couldn't fetch news"})
		return
	}

	c.JSON(http.StatusOK, models.NewsDetailsResponse{
		Message: "Fetched news successfully",
		News:    details,
	})
}
