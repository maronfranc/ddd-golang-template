package application

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/maronfranc/poc-golang-ddd/application/handler"
)

type Application struct{}

func (a *Application) ListenAndServe(port int) {
	r := chi.NewRouter()
	// r.Use(middleware.AllowContentType("application/json"))
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(applicationJson)
	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(middleware.Timeout(60 * time.Second))
	a.loadRoutes(r)

	log.Printf("Listening on port: %d", port)
	p := fmt.Sprintf(":%d", port)
	log.Fatal(http.ListenAndServe(p, r))
}

func (a *Application) loadRoutes(router chi.Router) {
	log.Print("APP loadRoutes")

	router.Mount("/examples", handler.RouteExample())
}

func applicationJson(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}