package main

import (
	"authServiceApp/data"
	"database/sql"
	"log"
	"os"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

const (
	gRpcPort = "50001"
)

type Config struct {
	Repo data.Repository
}

func main() {
	conn := connectToDB()
	if conn == nil {
		log.Panic("Can't connect to Postgres")
	}

	app := Config{}
	app.setupRepo(conn)

	app.gRPCListen()
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

var counts = 0

func connectToDB() *sql.DB {
	dsn := os.Getenv("DSN")

	for {
		// keep connecting to the database
		connection, err := openDB(dsn)
		if err != nil {
			log.Printf("Postgres is not yet ready")
			counts++
		} else {
			log.Printf("Connected to Postgres!")
			return connection
		}

		if counts > 10 {
			log.Println(err)
			return nil
		}

		log.Println("Backing off for two seconds...")
		time.Sleep(2 * time.Second)
	}
}

func (app *Config) setupRepo(conn *sql.DB) {
	db := data.NewRepository(conn)
	app.Repo = db
}
