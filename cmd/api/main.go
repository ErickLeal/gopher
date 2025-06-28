package main

import (
	"log"

	"github.com/ErickLeal/gopher/internal/db"
	"github.com/ErickLeal/gopher/internal/env"
	"github.com/ErickLeal/gopher/internal/store"
)

//	@title			GopherSocial API
//	@description	API for GopherSocial, a social network for gohpers
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

// @securityDefinitions.apikey	ApiKeyAuth
// @in							header
// @name						Authorization
// @description
func main() {
	env.LoadEnvs()
	cfg := config{
		addr:   env.SERVER_ADDR,
		apiUrl: env.API_URL,
		db: dbConfig{
			addr:         env.DB_ADDR,
			maxOpenConns: env.MAX_OPEN_CONNS,
			maxIdleConns: env.MAX_IDLE_CONNS,
			maxIdleTime:  env.MAX_IDLE_TIME,
		},
	}
	db, err := db.New(cfg.db.addr, cfg.db.maxOpenConns, cfg.db.maxIdleConns, cfg.db.maxIdleTime)
	if err != nil {
		log.Panic(err)
	}
	defer db.Close()
	log.Println("Connected to database")

	store := store.NewStorage(db)

	app := application{
		config: cfg,
		store:  store,
	}

	mux := app.mount()

	log.Fatal(app.run(mux))
}
