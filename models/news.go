package models

type News struct {
	Title         string `json:"title"`
	Description   string `json:"description"`
	Date          string `json:"date"`
	Category      string `json:"category"`
	Slug          string `json:"slug"`
	FeaturedImage string `json:"featured_image"`
}

type NewsResponse struct {
	Message string `json:"message"`
	News    []News `json:"news"`
}

type NewsDetailsContent struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

type NewsDetails struct {
	Title         string               `json:"tile"`
	FeaturedImage string               `json:"featured_image"`
	Date          string               `json:"date"`
	Source        string               `json:"source"`
	Content       []NewsDetailsContent `json:"content"`
}

type NewsDetailsResponse struct {
	Message string      `json:"message"`
	News    NewsDetails `json:"news"`
}
