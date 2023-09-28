package main

import (
	"log"

	"github.com/maronfranc/poc-golang-ddd/application"
)

func main() {
	log.Print("MAIN FILE")

	const PORT = 3000
	app := &application.Application{}
	app.ListenAndServe(PORT)
}
