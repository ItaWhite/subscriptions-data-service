package service

import (
	"context"
	"fmt"
	serviceErrors "subscriptions-data-service/internal/errors"
	"subscriptions-data-service/internal/model"
	"subscriptions-data-service/internal/repository"
	"time"
)

type RecordService struct {
	repo *repository.RecordRepository
}

func NewRecordService(r *repository.RecordRepository) *RecordService {
	return &RecordService{
		repo: r,
	}
}

func (s *RecordService) GetAll(ctx context.Context, limit, offset int) ([]model.Record, error) {
	if limit < 0 {
		return nil, fmt.Errorf("negative limit: %w", serviceErrors.ErrInvalidArgument)
	}
	if offset < 0 {
		return nil, fmt.Errorf("negative offset: %w", serviceErrors.ErrInvalidArgument)
	}
	return s.repo.GetAll(ctx, limit, offset)
}

func (s *RecordService) GetByID(ctx context.Context, id int) (model.Record, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *RecordService) Create(ctx context.Context, record model.Record) (model.Record, error) {
	if record.EndDate != nil && record.EndDate.Before(record.StartDate) {
		return model.Record{}, serviceErrors.ErrInvalidDates
	}
	return s.repo.Create(ctx, record)
}

func (s *RecordService) Update(ctx context.Context, id int, record model.Record) error {
	if record.EndDate != nil && record.EndDate.Before(record.StartDate) {
		return serviceErrors.ErrInvalidDates
	}
	return s.repo.Update(ctx, id, record)
}

func (s *RecordService) Delete(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}

func countMonth(from, to, start time.Time, end *time.Time) int {
	var endVal time.Time
	if end == nil {
		endVal = to
	} else {
		endVal = *end
	}
	if endVal.After(to) {
		endVal = to
	}
	if start.Before(from) {
		start = from
	}
	if start.After(endVal) {
		return 0
	}
	months := (endVal.Year()-start.Year())*12 + int(endVal.Month()-start.Month()) + 1
	return months

}

func (s *RecordService) TotalPrice(ctx context.Context, userIDStr string, serviceName string, from time.Time, to time.Time) (int, error) {
	if to.Before(from) {
		return 0, serviceErrors.ErrInvalidDates
	}
	pricesWithDates, err := s.repo.GetPricesInRange(ctx, userIDStr, serviceName, from, to)
	if err != nil {
		return 0, err
	}

	var total int
	for _, pwd := range pricesWithDates {
		price, start, end := pwd.Price, pwd.StartDate, pwd.EndDate
		total += price * countMonth(from, to, start, end)
	}

	return total, nil
}
