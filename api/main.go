package main

import (
	"fmt"
	"os"

	"github.com/maronfranc/poc-golang-ddd/application"
	"github.com/maronfranc/poc-golang-ddd/infrastructure"
	"github.com/maronfranc/poc-golang-ddd/infrastructure/database"
)

func main() {
	envfile, err := infrastructure.EnvGetFileName()
	if err != nil {
		panic(err)
	}
	err = infrastructure.EnvLoad(envfile)
	if err != nil {
		panic(err)
	}

	err = database.Start(envfile)
	if err != nil {
		panic(err)
	}
	defer database.CloseDb()

	port, err := infrastructure.EnvGet("API_PORT")
	if err != nil {
		panic(err)
	}

	app := &application.Application{}
	err = app.ListenAndServe(port)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Application listen failed: %v\n", err)
		panic(err)
	}
}
