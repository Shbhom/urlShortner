package main

import (
	"fmt"
	"github.com/go-playground/validator/v10"
)

type CreateUrlDTO struct {
	Url string `json:"url" validate:"https_url,required"`
}

func main() {
	v := validator.New()
	dto := CreateUrlDTO{Url: "http://google.com"}
	err := v.Struct(dto)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Success")
	}
}
