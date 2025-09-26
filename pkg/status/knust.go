package status

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

type StatusFunctions struct{}

func NewStatusFunctions() *StatusFunctions {
	return &StatusFunctions{}
}

func (s *StatusFunctions) GetKNUSTStatus() []models.KNUSTServer {
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

	return servers
}

func (sf *StatusFunctions) GetStatusBadge() ([]byte, error) {
	host := config.HostUrlProd

	if gin.Mode() == "debug" {
		host = config.HostUrlDev
	}

	url := fmt.Sprintf("%s/api/v1/knust-server-status", host)

	res, err := http.Get(url)
	if err != nil {
		return []byte{}, nil
	}

	var response models.KNUSTServerStatusResponse

	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		return []byte{}, err
	}

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
		return []byte{}, err
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return []byte{}, err
	}

	return body, nil
}
