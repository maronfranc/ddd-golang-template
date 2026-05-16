package example

import (
	"github.com/maronfranc/poc-golang-ddd/infrastructure/model"
)

type Service struct {
	ExampleRepository *Repository
}

func NewService(exampleRepository *Repository) *Service {
	return &Service{ExampleRepository: exampleRepository}
}

func (es *Service) GetMany(page, limit int) (*[]model.ManyExampleResponseDto, int) {
	return es.ExampleRepository.GetMany(page, limit)
}

func (es *Service) GetById(id string) (*model.Example, error) {
	return es.ExampleRepository.GetById(id)
}

func (es *Service) Create(e *model.CreateExampleDto) (*model.Example, error) {
	return es.ExampleRepository.Create(e)
}

func (es *Service) UpdateById(id string, e *model.CreateExampleDto) error {
	return es.ExampleRepository.UpdateById(id, e)
}

func (es *Service) DeleteById(id string) error {
	return es.ExampleRepository.DeleteById(id)
}
