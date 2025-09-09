package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"subscribe_service/internal/config"
	"subscribe_service/internal/entity"
	repository "subscribe_service/internal/repository/postgres"
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

// Create создаёт новую подписку
//
// @Summary      Создание новой подписки
// @Description  Создание новой подписки
// @Tags       	subscription
// @Accept       json
// @Produce		 json
// @Param        request  body  entity.CreateRequest  true  "DataForCreate"
// @Success      201  {object} entity.Subscription
// @Failure      400
// @Failure      500
// @Router       /subscription [post]
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

// GetRecordByID возвращает данные подписки по ID
//
// @Summary     Получение данных о подписке
// @Description Получение данных о подписке
// @Tags       	subscription
// @Produce		json
// @Param      	request  path  uuid.UUID  true  "Subscription ID"
// @Success    	200  {object} entity.Subscription
// @Failure    	400
// @Failure    	500
// @Router     	/subscription/{id} [get]
func (c *controller) GetRecordByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		log.Printf("Invalid ID format: %s. Error: %s\n", idStr, err)
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	resp, err := c.service.GetRecordByID(id)
	if err != nil {
		log.Printf("Failed to get record: %s. id: %s\n", err, id)
		http.Error(w, "Failed to get record", http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(resp)
	if err != nil {
		log.Println("marshalling error:", err)
		http.Error(w, "Failed to marshal response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	w.Write(data)
}

// Update обновляет данные подписки
//
// @Summary      Обновление данных подписки
// @Description  Обновление данных подписки
// @Tags       	subscription
// @Accept       json
// @Produce      json
// @Param        request  body  entity.Subscription  true  "Data for update"
// @Success      200  {object} entity.Subscription
// @Failure      400
// @Failure      500
// @Router       /subscription [put]
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

// Delete удаляет подписку по ID
//
// @Summary      Удаление подписки
// @Description  Удаление подписки по ID
// @Tags       	subscription
// @Param        id   path      string  true  "Subscription ID"
// @Success      204
// @Failure      400
// @Failure      500
// @Router       /subscription/{id} [delete]
func (c *controller) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		log.Printf("Invalid ID format: %s. Error: %s\n", idStr, err)
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	err = c.service.Delete(id)
	if err != nil {
		if errors.Is(err, repository.ErrIDNotFound) {
			log.Printf("%s %s\n", err, id)
			http.Error(w, fmt.Sprintf("Failed to delete record. %s: %s", err, id), http.StatusBadRequest)
			return
		}
		log.Printf("failed to delete record: %s\n", err)
		http.Error(w, "Failed to delete record", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// List возвращает список всех подписок
//
// @Summary      Получение списка подписок
// @Description  Получение списка всех подписок
// @Tags       	subscription
// @Produce      json
// @Success      200  {object} entity.SubscriptionsResponse
// @Failure      500
// @Router       /subscription [get]
func (c *controller) List(w http.ResponseWriter, r *http.Request) {
	resp, err := c.service.GetList()
	if err != nil {
		log.Printf("List error: %s", err)
		http.Error(w, "Getting list error.", http.StatusInternalServerError)
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

// GetTotal возвращает общую стоимость и количество подписок
//
// @Summary      Получение общей стоимости и количества подписок
// @Description  Получение общей стоимости и количества подписок с учётом фильтров.
// @Tags       	total
// @Accept       json
// @Produce      json
// @Param        request  body  entity.TotalRequest  true  "Фильтры для статистики. Все поля опциональны"
// @Success      200  {object} entity.TotalResponse
// @Failure      400
// @Failure      500
// @Router       /subscription/total [get]
func (c *controller) GetTotal(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("read request body error:", err)
		http.Error(w, "can't read body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var req entity.TotalRequest
	if err := json.Unmarshal(body, &req); err != nil {
		log.Println("unmarshalling error:", err)
		http.Error(w, "can't unmarshal body", http.StatusBadRequest)
		return
	}

	resp, err := c.service.GetTotal(&req)
	if err != nil {
		log.Println("getting total error:", err)
		http.Error(w, "can't get total", http.StatusInternalServerError)
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
