package main

import (
	"log"

	"github.com/ErickLeal/gopher/internal/env"
)

func main() {
	env.LoadEnvs()
	cfg := config{
		addr: env.SERVER_ADDR,
	}
	app := application{
		config: cfg,
	}

	mux := app.mount()

	log.Fatal(app.run(mux))
}
