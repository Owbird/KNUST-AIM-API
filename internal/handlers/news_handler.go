package handlers

import (
	"net/http"
	"strings"

	"github.com/Owbird/KNUST-AIM-Desktop-API/config"
	"github.com/Owbird/KNUST-AIM-Desktop-API/models"
	"github.com/anaskhan96/soup"
	"github.com/gin-gonic/gin"
)

func (h *Handlers) GetNewsHandler(c *gin.Context) {
	soup.Header("User-Agent", config.UserAgent)

	res, err := soup.Get(config.NewsEndpoint)

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Message: "Couldn't fetch news"})
	}

	htmlTree := soup.HTMLParse(res)

	newsRows := htmlTree.FindAll("div", "class", "views-row")

	appNews := []models.News{}

	for _, row := range newsRows {

		newsTag := row.Find("a")

		slug := strings.Split(newsTag.Attrs()["href"], "news-items")[1]

		titleTag := newsTag.Find("h3")

		descriptionTag := newsTag.Find("p")

		dateTag := newsTag.Find("span", "class", "post-date")

		date := strings.ReplaceAll(dateTag.Text(), "Published: ", "")

		categoryTag := newsTag.Find("span", "class", "post-cat")

		appNews = append(appNews, models.News{
			Title:       titleTag.Text(),
			Description: descriptionTag.Text(),
			Date:        strings.TrimSpace(date),
			Category:    strings.TrimSpace(categoryTag.Text()),
			Slug:        strings.TrimSpace(slug),
		})
	}

	c.JSON(http.StatusOK, models.NewsResponse{
		Message: "Fetched news successfully",
		News:    appNews,
	})

}
