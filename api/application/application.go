package application

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
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
	a.listenAndServe(r, port)
}
func (a *Application) listenAndServe(r *chi.Mux, port int) {
	p := fmt.Sprintf(":%d", port)
	srv := &http.Server{Addr: p, Handler: r}
	go func() {
		log.Printf("Listening on port: %d", port)
		err := srv.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("HTTP server error: %v", err)
		}
		log.Println("Stopped serving new connections.")
	}()
	interruptChan := make(chan os.Signal, 1)
	signal.Notify(interruptChan, syscall.SIGINT, syscall.SIGTERM)
	// Block until we receive our signal.
	<-interruptChan
	// Create a deadline to wait for.
	ctx, shutdownRelease := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownRelease()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("HTTP shutdown error: %v", err)
	}
	log.Println("Server graceful shutdown complete.")
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
	router.Use(a.applicationJson)
}
func (a *Application) LoadRoutes(router chi.Router) {
	router.Mount("/examples", handler.LoadExampleRoutes())
	router.NotFound(a.notFound)
}
func (a *Application) notFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	msg := fmt.Sprintf("Route not found(%s).", r.URL.Path)
	json.NewEncoder(w).Encode(&dto.ResponseError{Message: msg})
}
func (a *Application) applicationJson(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
