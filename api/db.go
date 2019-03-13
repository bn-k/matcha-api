package api

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
	"os"
)

func dbConnect() *sqlx.DB {
	connStr := "user=matcha password=secret dbname=matcha host=" +
		os.Getenv("POSTGRES_HOST") +
		" port=5432 sslmode=disable"
	db, err := sqlx.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	return db
}
