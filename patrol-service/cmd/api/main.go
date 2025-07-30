package main

import (
	"context"
	"log"
	"os"
	"patrolServiceApp/data"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	gRpcPort = "50002"
)

type Config struct {
	Repo     data.Repository
	FastRepo data.FastRepository
}

func main() {
	conn := connectToDB()
	if conn == nil {
		log.Panic("Can't connect to Postgres")
	}

	app := Config{}
	app.setupRepo(conn)
	app.setupFastRepo(app.Repo)

	app.gRPCListen()
}

func openDB(dsn string) (*pgxpool.Pool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// parse libpq-style DSN into config
	cfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}

	// create new pool with config
	cfg.MaxConns = 10
	cfg.MinConns = 1
	cfg.HealthCheckPeriod = time.Minute
	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, err
	}

	err = pool.Ping(ctx)
	if err != nil {
		return nil, err
	}

	return pool, nil
}

var counts = 0

func connectToDB() *pgxpool.Pool {
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

func (app *Config) setupRepo(conn *pgxpool.Pool) {
	db := data.NewRepository(conn)
	app.Repo = db
}

func (app *Config) setupFastRepo(persistRepo data.Repository) {
	db := data.NewFastRepository(persistRepo)
	app.FastRepo = db
}
