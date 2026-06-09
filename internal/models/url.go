package models

type CreateUrlDTO struct {
	Url string `json:"url" validate:"https_url,required"`
}
