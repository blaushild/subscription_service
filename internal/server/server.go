package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"subscribe_service/internal/config"
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
}

type Server struct {
	httpServer *http.Server
	sigCh      chan os.Signal
}

func NewServer(cfg *config.Config, c Controller) *Server {
	r := chi.NewRouter()
	setupRouter(r, c)

	return &Server{
		httpServer: &http.Server{
			Addr:         ":" + cfg.HttpServerPort,
			Handler:      r, // обработчик запросов
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
		},
	}
}

func (s *Server) Serve() {
	// Запуск сервера в отдельной горутине
	go func() {
		log.Println("Start subscribe service server...")
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("start subscribe service server error: %v", err)
		}
	}()

	// Ожидание сигнала остановки
	sig := <-s.sigCh
	log.Printf("Received signal: %s\n", sig)
}

// shutdownGracefully аккуратно останавливает сервер
func (s *Server) ShutdownGracefully() {
	// Создание контекста с таймаутом для graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Остановка сервера с использованием graceful shutdown
	if err := s.httpServer.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server stopped gracefully")
}
