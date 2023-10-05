package dto

type Response[T any] struct {
	Data T `json:"data"`
}

type ResponsePaginated[T any] struct {
	Data       *[]T       `json:"data"`
	Pagination *Paginated `json:"pagination"`
}

type ResponseError struct {
	Message string `json:"message"`
}

type Paginated struct {
	TotalRecord int    `json:"total_record"`
	TotalPage   int    `json:"total_page"`
	NextLink    string `json:"link_next,omitempty"`
	PrevLink    string `json:"link_prev,omitempty"`
}
