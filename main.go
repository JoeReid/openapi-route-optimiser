package main

import (
	"os"

	"github.com/JoeReid/openapi-route-optimiser/internal"
)

func main() {
	app, err := internal.New()
	if err != nil {
		os.Exit(1)
	}

	app.Run()
}
