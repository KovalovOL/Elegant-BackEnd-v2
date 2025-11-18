package db

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)


func getDBUrl() string{
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	url := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		dbUser, dbPassword, dbHost, dbPort, dbName,
	)
	return url
}

func GetPostgresDB() *pgx.Conn {
	url := getDBUrl()

	conn, err := pgx.Connect(context.Background(), url)
	if err != nil {
		panic("unable to connect to db")
	}
	return conn
}

func GetPool() *pgxpool.Pool {
	url := getDBUrl()
	pool, err := pgxpool.New(context.Background(), url)
	if err != nil {
		panic("unable to get db pool")
	}
	return pool
}
