package main

import (
	"fmt"
	"os"

	"github.com/maronfranc/poc-golang-ddd/application"
	"github.com/maronfranc/poc-golang-ddd/infrastructure"
	// _ "github.com/jackc/pgx/v5"
)

func main() {
	connStr := infrastructure.GetConnValues()
	err := infrastructure.ConnectDb(connStr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Connection failed: %v\n", err)
		panic(err)
	}

	const PORT = 3000
	app := &application.Application{}
	app.ListenAndServe(PORT)
	infrastructure.CloseDb()
}
