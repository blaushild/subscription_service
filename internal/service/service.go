package service

import (
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
}

func (s *service) GetRecordByID(id uuid.UUID) (*entity.Subscription, error) {
	return s.repo.GetRecordByID(id)
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
