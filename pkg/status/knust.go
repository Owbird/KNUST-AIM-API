package status

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/Owbird/KNUST-AIM-API/config"
	"github.com/Owbird/KNUST-AIM-API/models"
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

func (sf *StatusFunctions) GetStatusBadge() (string, error) {

	badge := "Up-green"

	for _, server := range sf.GetKNUSTStatus() {
		if server.Status == "Down" {
			badge = "Down-red"
			break
		}
	}

	shieldUrl := fmt.Sprintf("https://img.shields.io/badge/KNUST_Servers-%s", badge)

	return shieldUrl, nil
}
