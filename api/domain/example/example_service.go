package example

import (
	"github.com/maronfranc/poc-golang-ddd/application/dto"
	"github.com/maronfranc/poc-golang-ddd/infrastructure"
)

var exampleRepository = infrastructure.ExampleRepository{}

type ExampleService struct{}

func (es *ExampleService) GetMany(page, limit int) (*[]dto.ManyExampleResponseDto, int) {
	return exampleRepository.GetMany(page, limit)
}
func (es *ExampleService) GetById(id string) *dto.CreateExampleResponseDto {
	return exampleRepository.GetById(id)
}
func (es *ExampleService) Create(e *dto.CreateExampleDto) (*dto.CreateExampleResponseDto, error) {
	return exampleRepository.Create(e)
}
func (es *ExampleService) UpdateById(id string, e *dto.CreateExampleDto) error {
	return exampleRepository.UpdateById(id, e)
}
func (es *ExampleService) DeleteById(id string) error {
	return exampleRepository.DeleteById(id)
}
