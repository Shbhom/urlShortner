package main

import (
	"flag"

	"github.com/shbhom/urlShortner/internal/server"
)

func main() {
	env := flag.String("env", "local", "enviornment to run server")
	flag.Parse()

	server.Run(*env)
}
