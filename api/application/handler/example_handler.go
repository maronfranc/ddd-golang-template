package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/maronfranc/poc-golang-ddd/application/dto"
	"github.com/maronfranc/poc-golang-ddd/domain/example"
)

const ROUTE = "/examples"
const QS_PAGE = "page"
const QS_LIMIT = "limit"
const LIMIT = 10

func LoadExampleRoutes() chi.Router {
	r := chi.NewRouter()

	h := exampleHandler{}
	r.Get("/", h.getMany)
	r.Get("/{id}", h.getById)
	r.Post("/", h.create)
	r.Patch("/{id}", h.update)
	r.Delete("/{id}", h.deleteById)

	return r
}

var exampleService = &example.ExampleService{}

type exampleHandler struct{}

func (h *exampleHandler) getMany(w http.ResponseWriter, r *http.Request) {
	page := getQueryInt(r, QS_PAGE, 1)
	limit := getQueryInt(r, QS_LIMIT, 10)
	es, total := exampleService.GetMany(page, limit)
	pgn := &dto.ResponsePaginated[dto.ManyExampleResponseDto]{
		Data:       es,
		Pagination: buildPagination(ROUTE, total, page, limit)}
	json.NewEncoder(w).Encode(pgn)
}
func (h *exampleHandler) getById(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	e := exampleService.GetById(id)
	if e == nil {
		m := dto.ResponseError{Message: "Example not found"}
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(m)
		return
	}
	res := &dto.Response[dto.CreateExampleResponseDto]{Data: *e}
	json.NewEncoder(w).Encode(res)
}
func (h *exampleHandler) create(w http.ResponseWriter, r *http.Request) {
	var b dto.CreateExampleDto
	err := json.NewDecoder(r.Body).Decode(&b)
	if err != nil {
		m := dto.ResponseError{Message: "JSON decode error"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(m)
		return
	}
	c, err := exampleService.Create(&b)
	if err != nil {
		m := dto.ResponseError{Message: "Internal server error"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(m)
		return
	}
	res := &dto.Response[dto.CreateExampleResponseDto]{Data: *c}
	json.NewEncoder(w).Encode(res)
}
func (h *exampleHandler) update(w http.ResponseWriter, r *http.Request) {
	var b dto.CreateExampleDto
	err := json.NewDecoder(r.Body).Decode(&b)
	if err != nil {
		m := dto.ResponseError{Message: "JSON decode error"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(m)
		return
	}
	id := chi.URLParam(r, "id")
	err = exampleService.UpdateById(id, &b)
	res := &dto.Response[bool]{Data: err == nil}
	json.NewEncoder(w).Encode(res)
}
func (h *exampleHandler) deleteById(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	err := exampleService.DeleteById(id)
	res := &dto.Response[bool]{Data: err == nil}
	json.NewEncoder(w).Encode(res)
}

func getQueryInt(r *http.Request, n string, defaultV int) int {
	qv := r.URL.Query().Get(n)
	if qv == "" {
		return defaultV
	}
	v, err := strconv.Atoi(qv)
	if err != nil {
		v = defaultV
	}
	return v
}

func buildPagination(route string, total, page, limit int) *dto.Paginated {
	// SEE: https://stackoverflow.com/a/17974
	// SEE: https://stackoverflow.com/questions/17944/how-to-round-up-the-result-of-integer-division/17974
	totalPage := (total + limit - 1) / limit
	prevPage := page - 1
	nextPage := page + 1
	if prevPage < 0 {
		prevPage = 0
	}
	if nextPage > totalPage {
		nextPage = totalPage
	}
	prevLink := fmt.Sprintf("%s?page=%d&limit=%d", route, prevPage, limit)
	nextLink := fmt.Sprintf("%s?page=%d&limit=%d", route, nextPage, limit)
	return &dto.Paginated{
		TotalRecord: total,
		TotalPage:   int(totalPage),
		PrevLink:    prevLink,
		NextLink:    nextLink,
	}
}
