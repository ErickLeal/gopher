package main

import (
	"time"

	"github.com/ErickLeal/gopher/internal/db"
	"github.com/ErickLeal/gopher/internal/env"
	"github.com/ErickLeal/gopher/internal/mailer"
	"github.com/ErickLeal/gopher/internal/store"
	"go.uber.org/zap"
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
		mail: mailConfig{
			exp: time.Hour * 24 * 3, // 3 days
		},
	}

	// Logger
	logger := zap.Must(zap.NewProduction()).Sugar()
	defer logger.Sync()

	//Database connection
	db, err := db.New(cfg.db.addr, cfg.db.maxOpenConns, cfg.db.maxIdleConns, cfg.db.maxIdleTime)
	if err != nil {
		logger.Error(err)
	}
	defer db.Close()
	logger.Info("Connected to database")

	store := store.NewStorage(db)

	//mailer
	mailer := mailer.NewSendGridMailer(env.SENDGRID_API_KEY, env.FROM_EMAIL)

	app := application{
		config: cfg,
		store:  store,
		logger: logger,
		mailer: mailer,
	}

	mux := app.mount()

	logger.Fatal(app.run(mux))
}
