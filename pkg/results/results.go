package results

import (
	"fmt"
	"strings"

	"github.com/Owbird/KNUST-AIM-API/config"
	"github.com/Owbird/KNUST-AIM-API/internal/utils"
	"github.com/Owbird/KNUST-AIM-API/models"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
)

type ResultsFunctions struct{}

func NewResultsFunctions() *ResultsFunctions {
	return &ResultsFunctions{}
}

func (rf *ResultsFunctions) SelectResult(cookies any) (models.ResultsSelection, error) {
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
		return models.ResultsSelection{}, err
	}

	selectionUrl := fmt.Sprintf("%sResultChecker/AcademicSemSelection", config.BaseUrl)

	err = page.Navigate(selectionUrl)
	if err != nil {
		return models.ResultsSelection{}, err
	}

	err = page.WaitLoad()
	if err != nil {
		return models.ResultsSelection{}, err
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

	return models.ResultsSelection{
		Years: years,
		Sems:  sems,
	}, nil
}

func (rf *ResultsFunctions) GetResults(cookies any, payload models.GetResultsPayload) (models.GetResultsResponse, error) {
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
		return models.GetResultsResponse{}, err
	}

	selectionUrl := fmt.Sprintf("%sResultChecker/AcademicSemSelection", config.BaseUrl)

	err = page.Navigate(selectionUrl)
	if err != nil {
		return models.GetResultsResponse{}, err
	}

	err = page.WaitLoad()
	if err != nil {
		return models.GetResultsResponse{}, err
	}

	yearsEl := page.MustElement("select[id='ForminputState']")

	err = yearsEl.Select([]string{fmt.Sprintf("option[value='%s']", payload.Year)}, true, rod.SelectorTypeCSSSector)
	if err != nil {
		return models.GetResultsResponse{}, fmt.Errorf("Unknown year. Please use the available academic years.")
	}

	semsEl := page.MustElement("select[id='AcademicSemester']")

	err = semsEl.Select([]string{fmt.Sprintf("option[value='%s']", payload.Sem)}, true, rod.SelectorTypeCSSSector)
	if err != nil {
		return models.GetResultsResponse{}, fmt.Errorf("Unknown sem. Please use the available academic sems.")
	}

	forms := page.MustElements("form")

	displayResBtn := forms[1].MustElement("button[type='submit']")

	displayResBtn.MustClick()

	page.MustWaitNavigation()

	page.MustWaitLoad()

	sections := page.MustElements("table")

	if len(sections) == 0 {
		return models.GetResultsResponse{}, fmt.Errorf("Results not available for the chosen academic year and sem")
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
				return models.GetResultsResponse{}, err
			}

			text, err := el.Text()
			if err != nil {
				return models.GetResultsResponse{}, err
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

	return models.GetResultsResponse{
		Message: "Fetched results successfully",
		PersonalData: models.ResultsPersonalData{
			Name:      name,
			Year:      year,
			Sem:       payload.Sem,
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
	}, nil
}
