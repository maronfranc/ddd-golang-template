package handler

import (
	"encoding/json"
	"log"
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
	pgn := &dto.ResponsePaginated{Data: []interface{}{es}, Pagination: page}
	buf, err := json.Marshal(pgn)
	if err != nil {
		log.Println("JSON marshal error")
		return
	}
	w.Write(buf)
}

func (h *exampleHandler) getById(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	e := exampleService.GetById(id)
	if e == nil {
		http.Error(w, "Example not found", http.StatusNotFound)
		return
	}
	res := &dto.Response{Data: e}
	buf, err := json.Marshal(res)
	if err != nil {
		http.Error(w, "JSON marshal error", http.StatusInternalServerError)
		return
	}
	w.Write(buf)
}

func (h *exampleHandler) create(w http.ResponseWriter, r *http.Request) {
	log.Println("Handler create by id")
}

func (h *exampleHandler) update(w http.ResponseWriter, r *http.Request) {
	log.Println("Handler update by id")
}

func (h *exampleHandler) deleteById(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	b := exampleService.DeleteById(id)
	log.Println("Handler delete by id")
	log.Println(b)
}
