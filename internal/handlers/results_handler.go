package handlers

import (
	"fmt"
	"net/http"

	"github.com/Owbird/KNUST-AIM-API/config"
	"github.com/Owbird/KNUST-AIM-API/models"
	"github.com/gin-gonic/gin"
	"github.com/go-rod/rod/lib/proto"
)

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

	parsedCookies := cookies.(models.UserCookies)

	h.Browser.MustSetCookies(&proto.NetworkCookie{
		Name:     ".AspNetCore.Antiforgery.oBcnM5PKSJA",
		Value:    parsedCookies.Antiforgery,
		Path:     "/students",
		Domain:   "apps.knust.edu.gh",
		SameSite: "strict",
	}, &proto.NetworkCookie{
		Name:     ".AspNetCore.Identity.Application",
		Value:    parsedCookies.Identity,
		Path:     "/students",
		Domain:   "apps.knust.edu.gh",
		SameSite: "Lax",
	}, &proto.NetworkCookie{
		Name:     ".AspNetCore.Session",
		Value:    parsedCookies.Session,
		Path:     "/",
		Domain:   "apps.knust.edu.gh",
		SameSite: "Lax",
	})

	page := h.Browser.MustPage()

	defer page.Close()

	err := page.SetUserAgent(&proto.NetworkSetUserAgentOverride{
		UserAgent: config.UserAgent,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Message: "Couldn't get results. Please try again",
		})
	}

	selectionUrl := fmt.Sprintf("%sResultChecker/AcademicSemSelection", config.BaseUrl)

	err = page.Navigate(selectionUrl)

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Message: "Couldn't get results. Please try again",
		})
	}

	err = page.WaitLoad()

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Message: "Couldn't get results. Please try again",
		})
	}

	yearsEl := page.MustElement("select[id='ForminputState']")

	years := []string{}

	for _, child := range yearsEl.MustDescribe().Children {

		if child.Attributes[0] == "value" {
			years = append(years, child.Attributes[1])
		}

	}

	semsEl := page.MustElement("select[id='AcademicSemester']")

	sems := []string{}

	for _, child := range semsEl.MustDescribe().Children {

		if child.Attributes[0] == "value" {
			sems = append(sems, child.Attributes[1])
		}

	}

	c.JSON(http.StatusOK, models.ResultsSelectionResponse{
		Message: "Fetched results successfully",
		Results: models.ResultsSelection{
			Years: years,
			Sems:  sems,
		},
	})

}
