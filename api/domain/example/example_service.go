package example

import (
	"github.com/maronfranc/poc-golang-ddd/application/dto"
	"github.com/maronfranc/poc-golang-ddd/infrastructure"
)

var exampleRepository = infrastructure.ExampleRepository{}

type ExampleService struct{}

func (es *ExampleService) GetMany() (*[]dto.CreateExampleResponseDto, dto.Paginated) {
	examples, page := exampleRepository.GetMany()
	return examples, page
}
func (es *ExampleService) GetById(id string) *dto.CreateExampleResponseDto {
	return exampleRepository.GetById(id)
}
func (es *ExampleService) Create() *dto.CreateExampleResponseDto {
	return exampleRepository.Create()
}
func (es *ExampleService) Update() *dto.CreateExampleResponseDto {
	return exampleRepository.Update()
}
func (es *ExampleService) DeleteById(id string) bool {
	return exampleRepository.DeleteById(id)
}
