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
	NextLink    string `json:"next_page"`
	PrevLink    string `json:"prev_page"`
	// sort_by
	// sort_order 'asc' 'desc'
}
