package models

type CreateUrlDTO struct {
	Url string `json:"url" validate:"https_url,required"`
}

type UrlData struct {
	ShortCode string `json:"short_code"`
	TargetUrl string `json:"url"`
}
