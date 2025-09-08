package controller

import (
	"encoding/json"
	"io"
	"net/http"
	"subscribe_service/internal/config"
	"subscribe_service/internal/entity"
	"subscribe_service/internal/service"

	"log"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type SubscriptionService interface {
	Create(req *entity.CreateRequest) (*entity.Subscription, error)
	GetRecordByID(id uuid.UUID) (*entity.Subscription, error)
	Update(req *entity.Subscription) (*entity.Subscription, error)
	Delete(id uuid.UUID) error
	GetList() (*entity.SubscriptionsResponse, error)

	GetTotal(req *entity.TotalRequest) (*entity.TotalResponse, error)
}

type controller struct {
	service SubscriptionService
}

func NewController(cfg *config.Config) *controller {
	return &controller{
		service: service.NewService(cfg),
	}
}

func (c *controller) Create(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("read request body error:", err)
		http.Error(w, "can't read body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var req entity.CreateRequest
	if err := json.Unmarshal(body, &req); err != nil {
		log.Println("unmarshalling error:", err)
		http.Error(w, "can't unmarshal body", http.StatusBadRequest)
		return
	}

	resp, err := c.service.Create(&req)
	if err != nil {
		log.Println("create subscription error:", err)
		http.Error(w, "can't create subscription", http.StatusInternalServerError)
		return
	}

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Println("marshalling response error:", err)
		http.Error(w, "can't marshal response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	w.Write(jsonResp)
}

func (c *controller) GetRecordByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	resp, err := c.service.GetRecordByID(id)
	if err != nil {
		log.Printf("failed to get record: %s\n", err)
		http.Error(w, "Failed to get record", http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, "Failed to marshal response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	w.Write(data)
}

func (c *controller) Update(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("read request body error:", err)
		http.Error(w, "can't read body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var req entity.Subscription
	if err := json.Unmarshal(body, &req); err != nil {
		log.Println("unmarshalling error:", err)
		http.Error(w, "can't unmarshal body", http.StatusBadRequest)
		return
	}

	resp, err := c.service.Update(&req)
	if err != nil {
		log.Println("create subscription error:", err)
		http.Error(w, "can't create subscription", http.StatusInternalServerError)
		return
	}

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Println("marshalling response error:", err)
		http.Error(w, "can't marshal response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResp)
}

func (c *controller) Delete(w http.ResponseWriter, r *http.Request) {}

func (c *controller) List(w http.ResponseWriter, r *http.Request) {
	resp, err := c.service.GetList()
	if err != nil {
		log.Printf("List error: %s", err)
	}

	data, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, "Failed to marshal response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	w.Write(data)
}

func (c *controller) GetTotal(w http.ResponseWriter, r *http.Request) {}
