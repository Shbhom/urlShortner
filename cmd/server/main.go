package main

import (
	"os"

	"github.com/shbhom/urlShortner/internal/server"
)

func main() {

	env := os.Getenv("ENV")
	if env == "" {
		env = "local"
	}

	server.Run(env)
}
