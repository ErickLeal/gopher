package main

import (
	"log"

	"github.com/ErickLeal/gopher/internal/db"
	"github.com/ErickLeal/gopher/internal/env"
	"github.com/ErickLeal/gopher/internal/store"
)

func main() {
	env.LoadEnvs()

	conn, err := db.New(env.DB_ADDR, 3, 3, "15m")
	if err != nil {
		log.Panic(err)
	}
	defer conn.Close()
	log.Println("Connected to database")

	store := store.NewStorage(conn)
	db.Seed(store, conn)
}
