package models

type News struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Date        string `json:"date"`
	Category    string `json:"category"`
	Slug        string `json:"slug"`
}

type NewsResponse struct {
	Message string `json:"message"`
	News    []News `json:"news"`
}
