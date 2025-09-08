package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"subscribe_service/internal/config"
	"subscribe_service/internal/entity"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(cfg *config.Config) *Repository {
	dsn := fmt.Sprintf("host=%s dbname=%s user=%s password=%s sslmode=disable", cfg.DbHost, cfg.DbName, cfg.DbUser, cfg.DbPassword)

	db, err := connectDB(dsn, cfg.DbConnectAttempts)
	if err != nil {
		log.Fatal(err)
	}

	var schema = `
	CREATE TABLE IF NOT EXISTS subscriptions (
	    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		user_id UUID NOT NULL,
		start_date DATE NOT NULL,
		finish_date DATE,
	    service_name VARCHAR(255) NOT NULL,
	    price INTEGER NOT NULL CHECK (price >= 0),

		CONSTRAINT valid_date_range CHECK (finish_date IS NULL OR finish_date >= start_date)
	);`
	_, err = db.Exec(schema)
	if err != nil {
		log.Fatal("creating table error: ", err)
	}
	return &Repository{db: db}
}

func (r *Repository) Add(sub *entity.CreateRequest) (*entity.Subscription, error) {
	q := `INSERT INTO subscriptions (user_id, start_date, finish_date, service_name, price)
			VALUES (:UserID, :StartDate, :FinishDate, :ServiceName, :Price)
			RETURNING id, user_id, start_date, finish_date, service_name, price`

	params := map[string]any{
		"UserID":      sub.UserID,
		"StartDate":   sub.StartDate,
		"FinishDate":  sub.FinishDate,
		"ServiceName": sub.ServiceName,
		"Price":       sub.Price,
	}

	rows, err := r.db.NamedQuery(q, params)
	if err != nil {
		return nil, fmt.Errorf("db inserting error: %w", err)
	}
	defer rows.Close()

	var resp entity.Subscription

	if rows.Next() {
		err = rows.Scan(&resp.ID, &resp.UserID, &resp.StartDate, &resp.FinishDate, &resp.ServiceName, &resp.Price)
		if err != nil {
			return nil, fmt.Errorf("scanning id error: %w", err)
		}
	} else {
		return nil, fmt.Errorf("no rows returned from insert")
	}

	return &resp, nil
}

func (r *Repository) GetRecordByID(id uuid.UUID) (*entity.Subscription, error) {
	q := `SELECT id, user_id, start_date, finish_date, service_name, price
			FROM subscriptions
			WHERE id = $1`

	var resp entity.Subscription
	err := r.db.Get(&resp, q, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("subscription with id %s not found", id)
		}
		return nil, fmt.Errorf("getting of record %s error: %w", id, err)
	}

	return &resp, nil
}

func (r *Repository) Update(sub *entity.Subscription) (*entity.Subscription, error) {
	// TODO: проверить текущую запись и перенести данные, которые не изменяем
	q := `UPDATE subscriptions
		SET user_id = :UserID,
			start_date = :StartDate,
			finish_date = :FinishDate,
			service_name = :ServiceName,
			price = :Price
		WHERE id = :ID`

	params := map[string]any{
		"ID":          sub.ID,
		"UserID":      sub.UserID,
		"StartDate":   sub.StartDate,
		"FinishDate":  sub.FinishDate,
		"ServiceName": sub.ServiceName,
		"Price":       sub.Price,
	}

	_, err := r.db.NamedExec(q, params)
	if err != nil {
		return nil, fmt.Errorf("update subscription %s error: %w", sub.ID, err)
	}

	return sub, nil
}

func (r *Repository) Delete(id uuid.UUID) error {
	q := `DELETE FROM subscriptions WHERE id = $1`

	_, err := r.db.Exec(q, id)
	if err != nil {
		return fmt.Errorf("db Delete('%s') error: %w", id, err)
	}
	return nil
}

func (r *Repository) List() (*entity.SubscriptionsResponse, error) {
	q := `SELECT id, user_id, start_date, finish_date, service_name, price
			FROM subscriptions`
	rows, err := r.db.Query(q)
	if err != nil {
		return nil, fmt.Errorf("db List error: %w", err)
	}
	defer rows.Close()

	var resp entity.SubscriptionsResponse
	for rows.Next() {
		var sub entity.Subscription
		err := rows.Scan(&sub.ID, &sub.UserID, &sub.StartDate, &sub.FinishDate, &sub.ServiceName, &sub.Price)
		if err != nil {
			return nil, fmt.Errorf("row scan error: %w", err)
		}

		resp.Subscriptions = append(resp.Subscriptions, &sub)
	}

	return &resp, nil
}

func (r *Repository) GetTotal(req *entity.TotalRequest) (*entity.TotalResponse, error) {
	q := `SELECT SUM(price) AS total, COUNT(id) AS count 
			FROM subscriptions
			WHERE user_id = :UserID AND service_name = :ServiceName`

	rows, err := r.db.NamedQuery(q, req)
	if err != nil {
		return nil, fmt.Errorf("db List request error: %w", err)
	}
	defer rows.Close()

	var resp entity.TotalResponse
	if rows.Next() {
		err = rows.Scan(&resp.Total, &resp.Count)
		if err != nil {
			return nil, fmt.Errorf("scanning id error: %w", err)
		}
	} else {
		return nil, fmt.Errorf("no rows returned from insert")
	}

	return &resp, nil
}

func connectDB(dsn string, attempts int) (*sqlx.DB, error) {
	for i := 1; i <= attempts; i++ {
		db, err := sqlx.Connect("postgres", dsn)
		if err == nil {
			log.Printf("Connected to DB after %d attempts", i)
			return db, nil
		}
		time.Sleep(time.Duration(i) * time.Second)
	}

	return nil, fmt.Errorf("unable to connect to DB")
}
