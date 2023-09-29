package main

import (
	"github.com/maronfranc/poc-golang-ddd/application"
)

func main() {
	const PORT = 3000
	app := &application.Application{}
	app.ListenAndServe(PORT)
}
