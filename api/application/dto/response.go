package dto

type Response struct {
	Data any `json:"data"`
}

// type ResponsePaginated[T any] struct {
type ResponsePaginated struct {
	// Data []interface{} `json:"data"`
	// Data       []T       `json:"data"`
	Data       []any     `json:"data"`
	Pagination Paginated `json:"pagination"`
}

type Paginated struct {
	TotalRecords int `json:"total_records"`
	TotalPages   int `json:"total_pages"`
	CurrentPage  int `json:"current_page"`
	NextPage     int `json:"next_page"`
	PrevPage     int `json:"prev_page"`
	// sort_by
	// sort_order 'asc' 'desc'
}
