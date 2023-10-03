package application

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/maronfranc/poc-golang-ddd/application/dto"
	"github.com/maronfranc/poc-golang-ddd/application/handler"
)

type Application struct{}

func (a *Application) ListenAndServe(port int) {
	r := chi.NewRouter()
	a.LoadMiddlewares(r, true)
	a.LoadRoutes(r)

	log.Printf("Listening on port: %d", port)
	p := fmt.Sprintf(":%d", port)
	log.Fatal(http.ListenAndServe(p, r))
}

func (a *Application) LoadMiddlewares(router chi.Router, log bool) {
	if log {
		router.Use(middleware.RequestID)
		router.Use(middleware.RealIP)
		router.Use(middleware.Logger)
	}
	router.Use(middleware.Recoverer)
	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	router.Use(middleware.Timeout(60 * time.Second))
	router.Use(applicationJson)
}
func (a *Application) LoadRoutes(router chi.Router) {
	router.Mount("/examples", handler.LoadExampleRoutes())
	router.NotFound(a.notFound)
}
func (h *Application) notFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	msg := fmt.Sprintf("Route not found(%s).", r.URL.Path)
	json.NewEncoder(w).Encode(&dto.ResponseError{Message: msg})
}

func applicationJson(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
