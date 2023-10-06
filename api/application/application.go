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
	"github.com/maronfranc/poc-golang-ddd/infrastructure"
)

type Application struct{}

func (a *Application) ListenAndServe() error {
	r := a.Setup()
	port, err := infrastructure.EnvGet("API_PORT")
	if err != nil {
		return err
	}
	a.listenAndServe(r, port)
	return nil
}
func (a *Application) Setup() *chi.Mux {
	r := chi.NewRouter()
	a.LoadMiddlewares(r)
	a.LoadRoutes(r)
	return r
}
func (a *Application) listenAndServe(r *chi.Mux, port string) {
	p := fmt.Sprintf(":%s", port)
	srv := &http.Server{Addr: p, Handler: r}
	go func() {
		log.Printf("Listening on port: %s", port)
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
func (a *Application) LoadMiddlewares(router chi.Router) {
	log, _ := infrastructure.EnvGetAsBool("LOG")
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
