package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"subscribe_service/internal/config"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
)

type Controller interface {
	Create(w http.ResponseWriter, r *http.Request)
	GetRecordByID(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
	List(w http.ResponseWriter, r *http.Request)

	GetTotal(w http.ResponseWriter, r *http.Request)

	Swagger(w http.ResponseWriter, r *http.Request)
}

type Server struct {
	httpServer *http.Server
	sigCh      chan os.Signal
}

func NewServer(cfg *config.Config, c Controller) *Server {
	r := chi.NewRouter()
	setupRouter(r, c)

	s := &Server{
		httpServer: &http.Server{
			Addr:         ":" + cfg.HttpServerPort,
			Handler:      r, // обработчик запросов
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
		},
		sigCh: make(chan os.Signal, 1),
	}
	signal.Notify(s.sigCh, os.Interrupt, syscall.SIGTERM)
	return s
}

func (s *Server) Serve() {
	go func() {
		log.Println("Start subscribe service server...")
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("start subscribe service server error: %v", err)
		}
	}()

	sig := <-s.sigCh
	log.Printf("Received signal: %s\n", sig)
	s.ShutdownGracefully()
}

// shutdownGracefully аккуратно останавливает сервер
func (s *Server) ShutdownGracefully() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := s.httpServer.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server stopped gracefully")
}
