package repository

import (
	"context"
	"errors"
	"fmt"
	serviceErrors "subscriptions-data-service/internal/errors"
	"subscriptions-data-service/internal/model"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RecordRepository struct {
	db *pgxpool.Pool
}

func NewRecordRepository(pool *pgxpool.Pool) *RecordRepository {
	return &RecordRepository{
		db: pool,
	}
}

func (r *RecordRepository) GetAll(ctx context.Context) ([]model.Record, error) {
	rows, err := r.db.Query(ctx, "select * from records;")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var recordList []model.Record
	for rows.Next() {
		var record model.Record
		err = rows.Scan(&record.Id, &record.ServiceName, &record.Price, &record.UserID, &record.StartDate, &record.EndDate)
		if err != nil {
			return nil, err
		}
		recordList = append(recordList, record)
	}
	return recordList, nil
}

func (r *RecordRepository) GetByID(ctx context.Context, id int) (model.Record, error) {
	var record model.Record
	err := r.db.QueryRow(ctx, "select * from records where id=$1", id).
		Scan(&record.Id, &record.ServiceName, &record.Price, &record.UserID, &record.StartDate, &record.EndDate)
	if errors.Is(err, pgx.ErrNoRows) {
		return model.Record{}, serviceErrors.ErrRecordNotFound
	}
	if err != nil {
		return model.Record{}, err
	}
	return record, nil
}

func (r *RecordRepository) Create(ctx context.Context, record model.Record) (model.Record, error) {
	err := r.db.QueryRow(ctx, `
insert into records (service_name, price, user_id, start_date, end_date) 
values ($1, $2, $3, $4, $5) returning id;`, record.ServiceName, record.Price, record.UserID, record.StartDate, record.EndDate).Scan(&record.Id)
	if err != nil {
		return model.Record{}, err
	}
	return record, nil
}

func (r *RecordRepository) Update(ctx context.Context, id int, record model.Record) error {
	cmd, err := r.db.Exec(ctx, "update records set service_name=$1, price=$2, user_id=$3, start_date=$4, end_date=$5 where id=$6",
		record.ServiceName, record.Price, record.UserID, record.StartDate, record.EndDate, id)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return serviceErrors.ErrRecordNotFound
	}
	record.Id = id
	return err
}

func (r *RecordRepository) Delete(ctx context.Context, id int) error {
	cmd, err := r.db.Exec(ctx, "delete from records where id=$1", id)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return serviceErrors.ErrRecordNotFound
	}
	return nil
}

func (r *RecordRepository) GetPricesInRange(ctx context.Context, userIDStr string, serviceName string, from time.Time, to time.Time) ([]model.Record, error) {
	query := "select price, start_date, end_date from records"
	args := []any{from, to}
	query += " where (end_date is null or end_date >= $1) and start_date <= $2"
	if userIDStr != "" {
		args = append(args, userIDStr)
		query += fmt.Sprintf(" and user_id=$%d", len(args))
	}
	if serviceName != "" {
		args = append(args, serviceName)
		query += fmt.Sprintf(" and service_name=$%d", len(args))
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var pricesWithDates []model.Record
	for rows.Next() {
		var priceWithDates model.Record
		err = rows.Scan(&priceWithDates.Price, &priceWithDates.StartDate, &priceWithDates.EndDate)
		if err != nil {
			return nil, err
		}
		pricesWithDates = append(pricesWithDates, priceWithDates)
	}
	return pricesWithDates, nil
}
