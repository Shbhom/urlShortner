package services

import (
	"log"

	"github.com/sqids/sqids-go"
)

func NewSquid(minLength uint8) *sqids.Sqids {
	sq, err := sqids.New(sqids.Options{
		MinLength: minLength,
	})
	if err != nil {
		log.Fatal("Error while creating sqid Encoder")
	}
	return sq
}
