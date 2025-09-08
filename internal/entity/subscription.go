package entity

import "github.com/google/uuid"

type Subscription struct {
	ID          uuid.UUID   `json:"id" db:"id"`
	UserID      uuid.UUID   `json:"user_id" db:"user_id"`
	StartDate   CustomDate  `json:"start_date" db:"start_date"`
	FinishDate  *CustomDate `json:"finish_date" db:"finish_date"`
	ServiceName string      `json:"service_name" db:"service_name"`
	Price       uint        `json:"price" db:"price"`
}

type CreateRequest struct {
	UserID      uuid.UUID   `json:"user_id" db:"user_id"`
	StartDate   CustomDate  `json:"start_date" db:"start_date"`
	FinishDate  *CustomDate `json:"finish_date,omitempty" db:"finish_date"`
	ServiceName string      `json:"service_name" db:"service_name"`
	Price       uint        `json:"price" db:"price"`
}

type TotalRequest struct {
	UserID      uuid.UUID   `json:"user_id" db:"user_id"`
	ServiceName string      `json:"service_name" db:"service_name"`
	StartDate   *CustomDate `json:"start_date" db:"start_date"`
	FinishDate  *CustomDate `json:"finish_date,omitempty" db:"finish_date"`
}

type TotalResponse struct {
	UserID      uuid.UUID
	ServiceName string
	Total       uint
	Count       uint16
}

type SubscriptionsResponse struct {
	Subscriptions []*Subscription
}
