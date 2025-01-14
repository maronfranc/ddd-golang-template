package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/maronfranc/poc-golang-ddd/domain/dto"
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

var exampleService = &example.Service{}

type exampleHandler struct{}

func (h *exampleHandler) getMany(w http.ResponseWriter, r *http.Request) {
	page := getQueryInt(r, QS_PAGE, 1)
	limit := getQueryInt(r, QS_LIMIT, LIMIT)
	es, total := exampleService.GetMany(page, limit)
	pgn := &dto.ResponsePaginated[dto.ManyExampleResponseDto]{
		Data:       es,
		Pagination: buildPagination(ROUTE, total, page, limit)}
	json.NewEncoder(w).Encode(pgn)
}

func (h *exampleHandler) getById(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	e, err := exampleService.GetById(id)
	if errors.Is(err, sql.ErrNoRows) {
		encodeResponseError(w, "Example not found", http.StatusNotFound)
		return
	}
	if err != nil {
		encodeResponseError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	res := &dto.Response[dto.CreateExampleResponseDto]{Data: *e}
	json.NewEncoder(w).Encode(res)
}

func (h *exampleHandler) create(w http.ResponseWriter, r *http.Request) {
	var b dto.CreateExampleDto
	err := json.NewDecoder(r.Body).Decode(&b)
	if err != nil {
		encodeResponseError(w, "JSON decode error", http.StatusInternalServerError)
		return
	}
	c, err := exampleService.Create(&b)
	if err != nil {
		encodeResponseError(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	res := &dto.Response[dto.CreateExampleResponseDto]{Data: *c}
	json.NewEncoder(w).Encode(res)
}

func (h *exampleHandler) update(w http.ResponseWriter, r *http.Request) {
	var b dto.CreateExampleDto
	err := json.NewDecoder(r.Body).Decode(&b)
	if err != nil {
		encodeResponseError(w, "JSON decode error", http.StatusInternalServerError)
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

func getQueryInt(r *http.Request, name string, defaultV int) int {
	qv := r.URL.Query().Get(name)
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
	totalPage := int((total + limit - 1) / limit)
	prevPage := page - 1
	nextPage := page + 1
	var prevLink, nextLink string
	if prevPage > 1 {
		prevLink = fmt.Sprintf("%s?page=%d&limit=%d", route, prevPage, limit)
	}
	if nextPage < totalPage {
		nextLink = fmt.Sprintf("%s?page=%d&limit=%d", route, nextPage, limit)
	}
	return &dto.Paginated{
		TotalRecord: total,
		TotalPage:   totalPage,
		PrevLink:    prevLink,
		NextLink:    nextLink,
	}
}

func encodeResponseError(w http.ResponseWriter, msg string, code int) {
	m := dto.ResponseMessage{Message: msg}
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(m)
}
