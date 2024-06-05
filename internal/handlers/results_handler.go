package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/Owbird/KNUST-AIM-API/config"
	"github.com/Owbird/KNUST-AIM-API/internal/utils"
	"github.com/Owbird/KNUST-AIM-API/models"
	"github.com/gin-gonic/gin"
	"github.com/go-rod/rod"
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

	browser := utils.NewBrowser()

	browser.MustSetCookies(&proto.NetworkCookie{
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

	page := browser.MustPage()

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

	parsedCookies := cookies.(models.UserCookies)

	browser := utils.NewBrowser()

	browser.MustSetCookies(&proto.NetworkCookie{
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

	var resultsPayload models.GetResultsPayload

	err := c.BindJSON(&resultsPayload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Message: "Couldn't fetch user results. Please try again",
		})

		return
	}

	page := browser.MustPage()

	defer page.Close()

	err = page.SetUserAgent(&proto.NetworkSetUserAgentOverride{
		UserAgent: config.UserAgent,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Message: "Couldn't get results. Please try again",
		})

		return
	}

	selectionUrl := fmt.Sprintf("%sResultChecker/AcademicSemSelection", config.BaseUrl)

	err = page.Navigate(selectionUrl)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Message: "Couldn't get results. Please try again",
		})

		return
	}

	err = page.WaitLoad()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Message: "Couldn't get results. Please try again",
		})

		return
	}

	yearsEl := page.MustElement("select[id='ForminputState']")

	err = yearsEl.Select([]string{fmt.Sprintf("option[value='%s']", resultsPayload.Year)}, true, rod.SelectorTypeCSSSector)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Message: "Unknown year. Please use the available academic years.",
		})

		return
	}

	semsEl := page.MustElement("select[id='AcademicSemester']")

	err = semsEl.Select([]string{fmt.Sprintf("option[value='%s']", resultsPayload.Sem)}, true, rod.SelectorTypeCSSSector)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Message: "Unknown sem. Please use the available academic sems.",
		})

		return
	}

	forms := page.MustElements("form")

	displayResBtn := forms[1].MustElement("button[type='submit']")

	displayResBtn.MustClick()

	page.MustWaitNavigation()

	page.MustWaitLoad()

	sections := page.MustElements("table")

	if len(sections) == 0 {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Message: "Results not available for the chosen academic year and sem",
		})

		return
	}

	personalDataSection := sections[0]

	nameNYearEl := personalDataSection.MustElements("th")

	name := nameNYearEl[1].MustText()
	year := nameNYearEl[3].MustText()

	otherPersonalData := personalDataSection.MustElements("td")

	indexNo := otherPersonalData[0].MustText()
	programme := otherPersonalData[2].MustText()
	studentId := otherPersonalData[3].MustText()
	date := otherPersonalData[5].MustText()
	username := otherPersonalData[6].MustText()
	option := otherPersonalData[8].MustText()

	resultsSection := sections[1]

	tableRows := resultsSection.MustElements("tr")

	results := []models.Results{}

	for index, row := range tableRows {

		// Ignore headers
		if index == 0 {
			continue
		}

		result := models.Results{}

		for i := 0; i < 5; i++ {

			el, err := page.ElementFromNode(row.MustDescribe().Children[i])
			if err != nil {
				c.JSON(http.StatusInternalServerError, models.ErrorResponse{
					Message: "Couldn't get results. Please try again",
				})

				return
			}

			text, err := el.Text()
			if err != nil {
				c.JSON(http.StatusInternalServerError, models.ErrorResponse{
					Message: "Couldn't get results. Please try again",
				})

				return
			}

			switch i {
			case 0:
				result.CourseCode = text
			case 1:
				result.CourseName = text
			case 2:
				result.Credits = text
			case 3:
				result.TotalMark = text
			case 4:
				result.Grade = text
			}
		}

		results = append(results, result)

	}

	summarySection := sections[2].MustElements("tr")

	creditsRegisteredEl := summarySection[1].MustElements("td")

	creditsRegisteredSemester := creditsRegisteredEl[2].MustText()
	creditsRegisteredCumulative := creditsRegisteredEl[3].MustText()

	creditsObtainedEl := summarySection[2].MustElements("td")

	creditsObtainedSemester := creditsObtainedEl[1].MustText()
	creditsObtainedCumulative := creditsObtainedEl[2].MustText()

	creditsCalculatedEl := summarySection[3].MustElements("td")

	creditsCalculatedSemester := creditsCalculatedEl[1].MustText()
	creditsCalculatedCumulative := creditsCalculatedEl[2].MustText()

	weightedMarksEl := summarySection[4].MustElements("td")

	weightedMarksSemester := weightedMarksEl[1].MustText()
	weightedMarksCumulative := weightedMarksEl[2].MustText()

	cwaEl := summarySection[5].MustElements("td")

	cwaSemester := cwaEl[1].MustText()
	cwaCumulative := cwaEl[2].MustText()

	trails := strings.Split(summarySection[5].MustElement("th").MustText(), ", ")

	if trails[0] == "<none>" {
		trails = []string{}
	}

	c.JSON(http.StatusOK, models.GetResultsResponse{
		Message: "Fetched results successfully",
		PersonalData: models.ResultsPersonalData{
			Name:      name,
			Year:      year,
			IndexNo:   indexNo,
			Programme: programme,
			StudentID: studentId,
			Date:      date,
			Option:    option,
			Username:  username,
		},
		Results: results,
		Summary: models.ResultsSummary{
			CreditsRegistered: models.ResultsSummaryExtra{
				Semester:   creditsRegisteredSemester,
				Cumulative: creditsRegisteredCumulative,
			},
			CreditsObtained: models.ResultsSummaryExtra{
				Semester:   creditsObtainedSemester,
				Cumulative: creditsObtainedCumulative,
			},
			CreditsCalculated: models.ResultsSummaryExtra{
				Semester:   creditsCalculatedSemester,
				Cumulative: creditsCalculatedCumulative,
			},
			WeightedMarks: models.ResultsSummaryExtra{
				Semester:   weightedMarksSemester,
				Cumulative: weightedMarksCumulative,
			},
			CWA: models.ResultsSummaryExtra{
				Semester:   cwaSemester,
				Cumulative: cwaCumulative,
			},
		},
		Trails: trails,
	})
}
