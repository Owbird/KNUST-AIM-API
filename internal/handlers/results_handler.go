package handlers

import (
	"net/http"

	"github.com/Owbird/KNUST-AIM-API/models"
	"github.com/Owbird/KNUST-AIM-API/pkg/results"
	"github.com/gin-gonic/gin"
)

var resultsFunctions = results.NewResultsFunctions()

// @Summary Get available results
// @Description Returns a list of years and semester that the results are available for
// @Tags Results
// @Produce json
// @Success 200 {object} models.ResultsSelectionResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Router /user/results/selection [get]
func (h *Handlers) ResultSelectionHandler(c *gin.Context) {
	cookies, _ := c.Get("userCookies")

	res, err := resultsFunctions.SelectResult(cookies.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Message: "Couldn't get results. Please try again",
		})
	}

	c.JSON(http.StatusOK, models.ResultsSelectionResponse{
		Message: "Fetched results successfully",
		Results: res,
	})
}

// @Summary Get results
// @Description Returns results for the selected academic year and semester
// @Tags Results
// @Produce json
// @Accept  json
// @Param  year body string true "Year"
// @Param  sem body string true "Sem"
// @Success 200 {object} models.GetResultsResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Router /user/results [post]
func (h *Handlers) GetResultsHandler(c *gin.Context) {
	cookies, _ := c.Get("userCookies")

	var resultsPayload models.GetResultsPayload

	err := c.BindJSON(&resultsPayload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Message: "Couldn't fetch user results. Please try again",
		})

		return
	}

	res, err := resultsFunctions.GetResults(cookies.(string), resultsPayload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Message: "Couldn't fetch user results. Please try again",
		})
	}

	c.JSON(http.StatusOK, res)
}
