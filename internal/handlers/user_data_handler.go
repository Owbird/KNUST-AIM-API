package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/Owbird/KNUST-AIM-API/config"
	"github.com/Owbird/KNUST-AIM-API/models"
	"github.com/gin-gonic/gin"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
)

func (h *Handlers) GetUserData(c *gin.Context) {
	cookies, _ := c.Get("userCookies")

	parsedCookies := cookies.(models.UserCookies)

	var browser = rod.New().MustConnect().WithPanic(func(i interface{}) {
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
