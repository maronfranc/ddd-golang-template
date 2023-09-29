package infrastructure

import "github.com/maronfranc/poc-golang-ddd/application/dto"

func mockExample() *dto.CreateExampleResponseDto {
	return &dto.CreateExampleResponseDto{
		Id:          "aeiou-123",
		Title:       "Response title",
		Description: "Response descriptions",
	}
}

type ExampleRepository struct{}

func (er *ExampleRepository) GetMany() (*[]dto.CreateExampleResponseDto, dto.Paginated) {
	examples := &[]dto.CreateExampleResponseDto{*mockExample()}
	page := dto.Paginated{
		TotalRecords: 10,
		TotalPages:   10,
		CurrentPage:  1,
		NextPage:     2,
	}
	return examples, page
}
func (er *ExampleRepository) GetById(id string) *dto.CreateExampleResponseDto {
	return mockExample()
}
func (er *ExampleRepository) Create() *dto.CreateExampleResponseDto {
	return mockExample()
}
func (er *ExampleRepository) Update() *dto.CreateExampleResponseDto {
	return mockExample()
}
func (er *ExampleRepository) DeleteById(id string) bool {
	return true
}
