package db

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
)

func GetPostgresConnection() (*pgx.Conn){
	dbUser :=     os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost :=     os.Getenv("DB_HOST")
	dbPort :=     os.Getenv("DB_PORT")
	dbName :=     os.Getenv("DB_NAME")

	urlExample := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		dbUser, dbPassword, dbHost, dbPort, dbName,
	)

	conn, err := pgx.Connect(context.Background(), urlExample)
	if err != nil {
		panic("unable to connect to db")
	}
	return conn
}