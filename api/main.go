package main

import (
	"fmt"
	"os"

	"github.com/maronfranc/poc-golang-ddd/application"
	"github.com/maronfranc/poc-golang-ddd/infrastructure"
	// _ "github.com/jackc/pgx/v5"
)

func main() {
	envfile := infrastructure.EnvGetFile()
	err := infrastructure.Start(envfile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Infrastructure start failed: %v\n", err)
		panic(err)
	}

	app := &application.Application{}
	err = app.ListenAndServe()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Application listen failed: %v\n", err)
		panic(err)
	}

	infrastructure.CloseDb()
}
