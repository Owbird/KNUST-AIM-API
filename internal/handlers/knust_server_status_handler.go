package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"

	"github.com/Owbird/KNUST-AIM-API/config"
	"github.com/Owbird/KNUST-AIM-API/models"
	"github.com/gin-gonic/gin"
)

// @Summary Get the status of KNUST servers
// @Description This checks which of the used KNUST servers are up or down
// @Tags KNUST Servers
// @Produce json
// @Success 200 {object} models.KNUSTServerStatusResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /knust-server-status [get]
func (h *Handlers) KNUSTServerStatusHandler(c *gin.Context) {

	servers := []models.KNUSTServer{}

	wg := sync.WaitGroup{}

	urls := []string{config.MainUrl, config.AppsUrl}

	for _, url := range urls {

		wg.Add(1)
		go (func(url string) {

			defer wg.Done()

			mainUrlResp, err := http.Get(url)

			if err != nil {
				servers = append(servers, models.KNUSTServer{
					Url:    url,
					Status: "Down",
				})

				return
			}

			if mainUrlResp.StatusCode != http.StatusOK {
				servers = append(servers, models.KNUSTServer{
					Url:    url,
					Status: "Down",
				})

				return
			}

			servers = append(servers, models.KNUSTServer{
				Url:    url,
				Status: "Up",
			})

		})(url)
	}

	wg.Wait()

	c.JSON(http.StatusOK, models.KNUSTServerStatusResponse{
		Message: "Fetched server status successfully",
		Servers: servers,
	})

}

// @Summary Get the status of KNUST servers as a badge
// @Description This sums up the status of the servers and returns an SVG badge as a summary from shields.io
// @Tags KNUST Servers
// @Success 200 {string} string "OK"
// @Failure 500 {object} models.ErrorResponse
// @Router /knust-server-status/badge [get]
func (h *Handlers) KNUSTServerStatusBadgeHandler(c *gin.Context) {

	host := config.HostUrlProd

	if gin.Mode() == "debug" {
		host = config.HostUrlDev
	}

	url := fmt.Sprintf("%s/api/v1/knust-server-status", host)

	res, err := http.Get(url)

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Message: "Couldn't fetch badge"})
	}

	var response models.KNUSTServerStatusResponse

	json.NewDecoder(res.Body).Decode(&response)

	badge := "Up-green"

	for _, server := range response.Servers {
		if server.Status == "Down" {
			badge = "Down-red"
			break
		}
	}

	shieldUrl := fmt.Sprintf("https://img.shields.io/badge/KNUST_Servers-%s", badge)

	res, err = http.Get(shieldUrl)

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Message: "Couldn't fetch badge"})
	}

	body, err := io.ReadAll(res.Body)

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Message: "Couldn't fetch badge"})
	}

	c.Data(http.StatusOK, "image/svg+xml;charset=utf-8", body)

}
