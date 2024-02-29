package handlers

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/Owbird/KNUST-AIM-Desktop-API/config"
	"github.com/Owbird/KNUST-AIM-Desktop-API/models"
	"github.com/gin-gonic/gin"
)

func (h *Handlers) KNUSTServerStatusHandler(c *gin.Context) {

	servers := []models.KNUSTServer{}

	wg := sync.WaitGroup{}

	urls := []string{config.MainUrl, config.AppsUrl}

	badge := "Up-green"

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

				badge = "Down-red"

				return
			}

			if mainUrlResp.StatusCode != http.StatusOK {
				servers = append(servers, models.KNUSTServer{
					Url:    url,
					Status: "Down",
				})

				badge = "Down-red"

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
		Badge:   fmt.Sprintf("https://img.shields.io/badge/KNUST_Servers-%s", badge),
	})

}
