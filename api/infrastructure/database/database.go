package database

import (
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/maronfranc/poc-golang-ddd/infrastructure"
)

var DbConn *sqlx.DB

func Start(envfile string) error {
	connStr, err := infrastructure.EnvGet("PG_CONNECTION_STR")
	if err != nil {
		return err
	}

	err = connectDb(connStr)
	if err != nil {
		return err
	}

	err = DbConn.Ping()
	if err != nil {
		return err
	}

	return err
}

func connectDb(connectionString string) error {
	conn, err := sqlx.Open("postgres", connectionString)
	DbConn = conn
	return err
}

func CloseDb() {
	log.Print("Closing infrastructure connection")
	DbConn.Close()
}
