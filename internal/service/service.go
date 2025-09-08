package service

import (
	"fmt"
	"subscribe_service/internal/config"
	"subscribe_service/internal/entity"
	repository "subscribe_service/internal/repository/postgres"

	"github.com/google/uuid"
)

type Repository interface {
	Add(sub *entity.CreateRequest) (*entity.Subscription, error)
	GetRecordByID(id uuid.UUID) (*entity.Subscription, error)
	Update(sub *entity.Subscription) (*entity.Subscription, error)
	Delete(id uuid.UUID) error
	List() (*entity.SubscriptionsResponse, error)
	GetTotal(req *entity.TotalRequest) (*entity.TotalResponse, error)
}

type service struct {
	repo Repository
}

func NewService(cfg *config.Config) *service {
	return &service{
		repo: repository.NewRepository(cfg),
	}
}
func (s *service) Create(req *entity.CreateRequest) (*entity.Subscription, error) {
	return s.repo.Add(req)
	// if err != nil {
	// 	return nil, err
	// }

	// return &entity.Subscription{
	// 	ID:          resp.ID,
	// 	UserID:      req.UserID,
	// 	StartDate:   req.StartDate,
	// 	FinishDate:  req.FinishDate,
	// 	ServiceName: req.ServiceName,
	// 	Price:       req.Price,
	// }, nil
}

func (s *service) GetRecordByID(id uuid.UUID) (*entity.Subscription, error) {

	resp, err := s.repo.GetRecordByID(id)
	if err != nil {
		return nil, fmt.Errorf("GetRecordByID error: %w", err)
	}

	return resp, nil
}

func (s *service) Update(req *entity.Subscription) (*entity.Subscription, error) {
	return s.repo.Update(req)
}

func (s *service) Delete(id uuid.UUID) error {
	return s.repo.Delete(id)
}

func (s *service) GetList() (*entity.SubscriptionsResponse, error) {
	return s.repo.List()
}

func (s *service) GetTotal(req *entity.TotalRequest) (*entity.TotalResponse, error) {
	return s.repo.GetTotal(req)
}

//
// 	Update(req *entity.Subscription) (*entity.Subscription, error)
// 	Delete(id uuid.UUID) error
// 	GetList() (*entity.SubscriptionsResponse, error)

// 	GetTotal(req entity.TotalRequst) (entity.TotalResponse, error)
