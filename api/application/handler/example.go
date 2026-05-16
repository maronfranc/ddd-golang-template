package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/maronfranc/poc-golang-ddd/domain/dto"
	"github.com/maronfranc/poc-golang-ddd/domain/example"
	"github.com/maronfranc/poc-golang-ddd/infrastructure/model"
	"github.com/maronfranc/poc-golang-ddd/util"
)

func LoadExampleRoutes(
	route string,
	exampleService *example.Service,
) chi.Router {
	r := chi.NewRouter()

	h := NewExampleHandler(route, exampleService)
	r.Get("/", h.getMany)
	r.Get("/{id}", h.getById)
	r.Post("/", h.create)
	r.Patch("/{id}", h.update)
	r.Delete("/{id}", h.deleteById)

	return r
}

type exampleHandler struct {
	route          string
	exampleService *example.Service
}

func NewExampleHandler(
	route string,
	exampleService *example.Service,
) *exampleHandler {
	return &exampleHandler{
		route:          route,
		exampleService: exampleService,
	}
}

func (eh *exampleHandler) getMany(w http.ResponseWriter, r *http.Request) {
	page := util.GetQueryInt(r, util.QS_PAGE, 1)
	limit := util.GetQueryInt(r, util.QS_LIMIT, util.LIMIT)
	es, total := eh.exampleService.GetMany(page, limit)
	pgn := &dto.ResponsePaginated[model.ManyExampleResponseDto]{
		Data:       es,
		Pagination: util.BuildPagination(eh.route, total, page, limit),
	}
	json.NewEncoder(w).Encode(pgn)
}

func (eh *exampleHandler) getById(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	e, err := eh.exampleService.GetById(id)
	if errors.Is(err, sql.ErrNoRows) {
		util.EncodeResponseError(w, "Example not found", http.StatusNotFound)
		return
	}
	if err != nil {
		util.EncodeResponseError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res := &dto.Response[model.Example]{Data: *e}
	json.NewEncoder(w).Encode(res)
}

func (eh *exampleHandler) create(w http.ResponseWriter, r *http.Request) {
	var b model.CreateExampleDto
	err := json.NewDecoder(r.Body).Decode(&b)
	if err != nil {
		util.EncodeResponseError(w, "JSON decode error", http.StatusInternalServerError)
		return
	}

	c, err := eh.exampleService.Create(&b)
	if err != nil {
		util.EncodeResponseError(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	res := &dto.Response[model.Example]{Data: *c}
	json.NewEncoder(w).Encode(res)
}

func (eh *exampleHandler) update(w http.ResponseWriter, r *http.Request) {
	var b model.CreateExampleDto
	err := json.NewDecoder(r.Body).Decode(&b)
	if err != nil {
		util.EncodeResponseError(w, "JSON decode error", http.StatusInternalServerError)
		return
	}

	id := chi.URLParam(r, "id")
	err = eh.exampleService.UpdateById(id, &b)
	res := &dto.Response[bool]{Data: err == nil}
	json.NewEncoder(w).Encode(res)
}

func (eh *exampleHandler) deleteById(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	err := eh.exampleService.DeleteById(id)
	res := &dto.Response[bool]{Data: err == nil}
	json.NewEncoder(w).Encode(res)
}
