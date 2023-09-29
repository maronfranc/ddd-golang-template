package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/maronfranc/poc-golang-ddd/application/dto"
	"github.com/maronfranc/poc-golang-ddd/domain/example"
)

func RouteExample() chi.Router {
	r := chi.NewRouter()

	h := exampleHandler{}
	r.Get("/", h.getMany)
	r.Get("/{id}", h.getById)
	r.Post("/", h.create)
	r.Patch("/", h.update)
	r.Delete("/{id}", h.deleteById)

	return r
}

var exampleService = &example.ExampleService{}

type exampleHandler struct{}

func (h *exampleHandler) getMany(w http.ResponseWriter, r *http.Request) {
	es, page := exampleService.GetMany()
	pgn := &dto.ResponsePaginated[dto.CreateExampleResponseDto]{Data: es, Pagination: page}
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
	res := &dto.Response{Data: e}
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
	c := exampleService.Create(b)
	res := &dto.Response{Data: c}
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
	u := exampleService.Update(b)
	res := &dto.Response{Data: u}
	json.NewEncoder(w).Encode(res)
}
func (h *exampleHandler) deleteById(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	d := exampleService.DeleteById(id)
	res := &dto.Response{Data: d}
	json.NewEncoder(w).Encode(res)
}
