package entity

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

type CustomDate time.Time

const customDateFormat = "01-2006"

func (cd *CustomDate) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	if s == "null" || s == "" {
		*cd = CustomDate(time.Time{})
		return nil
	}
	t, err := time.Parse(customDateFormat, s)
	if err != nil {
		// Fallback to standard time format
		t, err = time.Parse(time.RFC3339, s)
		if err != nil {
			// Fallback to date-only format
			t, err = time.Parse("2006-01-02", s)
			if err != nil {
				return fmt.Errorf("failed to parse time in any supported format: %s", s)
			}
		}
	}
	*cd = CustomDate(t)
	return nil
}

func (cd CustomDate) MarshalJSON() ([]byte, error) {
	if time.Time(cd).IsZero() {
		return json.Marshal(nil)
	}
	return json.Marshal(time.Time(cd).Format(customDateFormat))
}

func (cd *CustomDate) Scan(value interface{}) error {
	if value == nil {
		*cd = CustomDate(time.Time{})
		return nil
	}
	switch v := value.(type) {
	case time.Time:
		*cd = CustomDate(v)
		return nil
	}
	return fmt.Errorf("can't convert %T to CustomDate", value)
}

func (cd CustomDate) Value() (driver.Value, error) {
	if time.Time(cd).IsZero() {
		return nil, nil
	}
	return time.Time(cd), nil
}

type Subscription struct {
	ID          uuid.UUID   `json:"id" db:"id"`
	UserID      uuid.UUID   `json:"user_id" db:"user_id"`
	StartDate   CustomDate  `json:"start_date" db:"start_date"`
	FinishDate  *CustomDate `json:"finish_date" db:"finish_date"`
	ServiceName string      `json:"service_name" db:"service_name"`
	Price       uint        `json:"price" db:"price"`
}

type SubscriptionRequest struct {
	ID uuid.UUID
}

type CreateRequest struct {
	UserID      uuid.UUID   `json:"user_id" db:"user_id"`
	StartDate   CustomDate  `json:"start_date" db:"start_date"`
	FinishDate  *CustomDate `json:"finish_date,omitempty" db:"finish_date"`
	ServiceName string      `json:"service_name" db:"service_name"`
	Price       uint        `json:"price" db:"price"`
}

type TotalRequest struct {
	UserID      uuid.UUID
	ServiceName string
	StartDate   *CustomDate
	FinishDate  *CustomDate
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
