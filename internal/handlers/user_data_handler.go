package handlers

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/Owbird/KNUST-AIM-API/config"
	"github.com/Owbird/KNUST-AIM-API/models"
	"github.com/gin-gonic/gin"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"
)

// @Summary Get User Data
// @Description Returns personal, programme and contact user data
// @Tags User
// @Produce json
// @Success 200 {object} models.UserDataResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Router /user [get]
func (h *Handlers) GetUserData(c *gin.Context) {
	cookies, _ := c.Get("userCookies")

	parsedCookies := cookies.(models.UserCookies)

	controlUrl := launcher.New().NoSandbox(true).MustLaunch()

	var browser = rod.New().ControlURL(controlUrl).MustConnect().WithPanic(func(i interface{}) {
		log.Println("[!] Headerless browser proberly lost context.")
	})

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

	page.SetUserAgent(&proto.NetworkSetUserAgentOverride{
		UserAgent: config.UserAgent,
	})

	profileUrl := fmt.Sprintf("%sHome/StudentProfile", config.BaseUrl)

	page.Navigate(profileUrl)

	page.WaitLoad()

	schoolEmail := *page.MustElement("input[name='StduentDTO.SchoolEmail']").MustAttribute("value")

	userName := strings.Split(schoolEmail, "@")[0]

	surname := *page.MustElement("input[name='StduentDTO.Surname']").MustAttribute("value")
	otherNames := *page.MustElement("input[name='StduentDTO.OtherName']").MustAttribute("value")
	studentId := *page.MustElement("input[name='StduentDTO.StudentId']").MustAttribute("value")

	indexNo, _ := page.MustElements("span[style='font-weight:600;']")[1].Text()

	dateOfBirth := *page.MustElement("input[name='date']").MustAttribute("value")
	gender := *page.MustElement("input[name='StduentDTO.Gender']").MustAttribute("value")
	knustMobile := *page.MustElement("input[name='StduentDTO.SchoolMobile']").MustAttribute("value")
	personalEmail := *page.MustElement("input[name='StduentDTO.OtherEmail']").MustAttribute("value")
	primaryMobile := *page.MustElement("input[name='StduentDTO.PrimaryMobile']").MustAttribute("value")
	country := *page.MustElement("input[name='StduentDTO.Country']").MustAttribute("value")
	otherPhone := *page.MustElement("input[name='StduentDTO.OtherPhone']").MustAttribute("value")

	resAddress1 := *page.MustElement("input[name='StduentDTO.ResAdd1']").MustAttribute("value")
	resAddress2 := *page.MustElement("input[name='StduentDTO.ResAdd2']").MustAttribute("value")
	resAddress3 := *page.MustElement("input[name='StduentDTO.ResAdd3']").MustAttribute("value")
	resAddress4 := *page.MustElement("input[name='StduentDTO.ResAdd4']").MustAttribute("value")

	postAddress1 := *page.MustElement("input[name='StduentDTO.PostAdd1']").MustAttribute("value")
	postAddress2 := *page.MustElement("input[name='StduentDTO.PostAdd2']").MustAttribute("value")
	postAddress3 := *page.MustElement("input[name='StduentDTO.PostAdd3']").MustAttribute("value")
	postAddress4 := *page.MustElement("input[name='StduentDTO.PostAdd4']").MustAttribute("value")

	regionNReligionEl := page.MustElements("option[selected='selected']")

	region, _ := regionNReligionEl[0].Text()
	religion, _ := regionNReligionEl[1].Text()

	programmeStream, _ := page.MustElements("h5")[1].Text()

	userData := models.UserData{
		Personal: models.PersonalUserData{
			Username:    userName,
			Surname:     surname,
			OtherNames:  otherNames,
			Gender:      gender,
			DateOfBirth: dateOfBirth,
			Country:     country,
			Region:      region,
			Religion:    religion,
		},
		Programme: models.ProgrammeUserData{
			StudentId:       studentId,
			IndexNo:         indexNo,
			ProgrammeStream: programmeStream,
		},
		Contact: models.ContactUserData{
			SchoolEmail:        schoolEmail,
			PersonalEmail:      personalEmail,
			KNUSTMobile:        knustMobile,
			PersonalMobile:     primaryMobile,
			AltPersonalMobile:  otherPhone,
			ResidentialAddress: fmt.Sprintf("%s\n%s\n%s\n%s\n", resAddress1, resAddress2, resAddress3, resAddress4),
			PostalAddress:      fmt.Sprintf("%s\n%s\n%s\n%s\n", postAddress1, postAddress2, postAddress3, postAddress4),
		},
	}

	c.JSON(http.StatusOK, models.UserDataResponse{
		Message:  "Fetched user data successfully",
		UserData: userData,
	})

}

// @Summary User image
// @Description Serves up the user image based on the student id
// @Tags User
// @Produce json
// @Param  studentId path string true "Student ID"
// @Success 200 {string} string "OK"
// @Failure 500 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 400 {object} models.ErrorResponse
// @Router /user/image/{studentId} [get]
func (h *Handlers) GetUserImage(c *gin.Context) {

	id, ok := c.Params.Get("id")

	if !ok {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Message: "Invalid user id",
		})

		return
	}

	url := fmt.Sprintf("%s?id=%s", config.UserImageUrl, id)

	resp, err := http.Get(url)

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Message: "Couldn't fetch user image",
		})

		return
	}

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Message: "Couldn't fetch user image",
		})

		return
	}

	c.Data(http.StatusOK, "image/jpg", body)
}
