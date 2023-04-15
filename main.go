package main

import (
	"database/sql"
	"log"

	"github.com/Munchies-Engineering/backend/api"
	db "github.com/Munchies-Engineering/backend/db/sqlc"
	"github.com/Munchies-Engineering/backend/util"
	_ "github.com/lib/pq"
)

func main() {
	config, err := util.LoadConfig("./")
	if err != nil {
		log.Fatal(err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	store := db.NewStore(conn)
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal(err)
	}
	server.Start(config.ServerAddress)
}
