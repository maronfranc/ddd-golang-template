package example

import (
	"log"
)

type ExampleService struct{}

func (es *ExampleService) GetMany() string {
	log.Println("Example GET MANY")
	return "TODO example getmany"
}

func (es *ExampleService) GetById() {
	log.Println("Example GET By ID")
}
