package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/Owbird/KNUST-AIM-API/config"
	"github.com/Owbird/KNUST-AIM-API/models"
	"github.com/anaskhan96/soup"
	"github.com/gin-gonic/gin"
)

// @Summary Get latest news
// @Description Returns the latest news available
// @Tags News
// @Produce json
// @Success 200 {object} models.NewsResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /news [get]
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

		slug = strings.ReplaceAll(slug, "/", "")

		titleTag := newsTag.Find("h3")

		descriptionTag := newsTag.Find("p")

		dateTag := newsTag.Find("span", "class", "post-date")

		date := strings.ReplaceAll(dateTag.Text(), "Published: ", "")

		categoryTag := newsTag.Find("span", "class", "post-cat")

		featuredImage := newsTag.Find("img").Attrs()["src"]

		appNews = append(appNews, models.News{
			Title:         titleTag.Text(),
			Description:   descriptionTag.Text(),
			Date:          strings.TrimSpace(date),
			Category:      strings.TrimSpace(categoryTag.Text()),
			Slug:          strings.TrimSpace(slug),
			FeaturedImage: fmt.Sprintf("%s%s", config.MainUrl, featuredImage),
		})
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
	soup.Header("User-Agent", config.UserAgent)

	slug := c.Param("slug")

	newsEndpoint := fmt.Sprintf("%s/news-items/%s", config.NewsEndpoint, slug)

	res, err := soup.Get(newsEndpoint)

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Message: "Couldn't fetch news"})
	}

	htmlTree := soup.HTMLParse(res)

	newsInfoTag := htmlTree.Find("div", "class", "post-info")

	title := newsInfoTag.Find("h2")

	featuredImageTag := newsInfoTag.Find("div", "class", "featured-img").Find("img")

	postDateTag := newsInfoTag.Find("span", "class", "post-date")

	date := strings.ReplaceAll(postDateTag.Text(), "Published: ", "")

	postSourceTag := newsInfoTag.Find("span", "class", "post-source")

	source := strings.ReplaceAll(postSourceTag.Text(), "Source: ", "")

	articleContentTag := newsInfoTag.Find("div", "class", "article-content").Find("div")

	content := []models.NewsDetailsContent{}

	for _, child := range articleContentTag.Children() {

		switch child.Pointer.Data {
		case "p":
			content = append(content, models.NewsDetailsContent{
				Type:  "text",
				Value: child.Text(),
			})

		case "figure":
			img := child.Find("img").Attrs()["src"]

			content = append(content, models.NewsDetailsContent{
				Type:  "media",
				Value: fmt.Sprintf("%s%s", config.MainUrl, img),
			})
		}
	}

	c.JSON(http.StatusOK, models.NewsDetailsResponse{
		Message: "Fetched news successfully",
		News: models.NewsDetails{
			Title:         title.Text(),
			FeaturedImage: fmt.Sprintf("%s%s", config.MainUrl, featuredImageTag.Attrs()["src"]),
			Date:          strings.TrimSpace(date),
			Source:        strings.TrimSpace(source),
			Content:       content,
		},
	})

}
