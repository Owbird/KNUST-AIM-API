package models

type KNUSTServer struct {
	Status string `json:"status"`
	Url    string `json:"url"`
}

type KNUSTServerStatusResponse struct {
	Message string        `json:"message"`
	Servers []KNUSTServer `json:"servers"`
}
