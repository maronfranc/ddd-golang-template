package handler

import (
	"log"
	"net/http"

	"github.com/go-chi/chi"
)

func RouteExample() chi.Router {
	r := chi.NewRouter()
	log.Println("Route Example")

	h := exampleHandler{}
	r.Get("/", h.getMany)
	r.Get("/{id}", h.getById)
	r.Post("/", h.create)
	r.Patch("/", h.update)
	r.Delete("/{id}", h.deleteById)

	return r
}

type exampleHandler struct{}

func (h *exampleHandler) getMany(w http.ResponseWriter, r *http.Request) {
	log.Println("Handler get many")
}

func (h *exampleHandler) getById(w http.ResponseWriter, r *http.Request) {
	log.Println("Handler get by id")
}

func (h *exampleHandler) deleteById(w http.ResponseWriter, r *http.Request) {
	log.Println("Handler delete by id")
}

func (h *exampleHandler) create(w http.ResponseWriter, r *http.Request) {
	log.Println("Handler create by id")
}

func (h *exampleHandler) update(w http.ResponseWriter, r *http.Request) {
	log.Println("Handler update by id")
}
