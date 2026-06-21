package internal

import (
	"context"
	"time"
)

type recordService struct {
	repo *recordRepository
}

func NewRecordService(r *recordRepository) *recordService {
	return &recordService{
		repo: r,
	}
}

func (s *recordService) GetAll(ctx context.Context) ([]Record, error) {
	return s.repo.GetAll(ctx)
}

func (s *recordService) GetByID(ctx context.Context, id int) (Record, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *recordService) Create(ctx context.Context, record Record) (Record, error) {
	if record.EndDate != nil && record.EndDate.Before(record.StartDate) {
		return Record{}, ErrInvalidDates
	}
	return s.repo.Create(ctx, record)
}

func (s *recordService) Update(ctx context.Context, id int, record Record) error {
	if record.EndDate != nil && record.EndDate.Before(record.StartDate) {
		return ErrInvalidDates
	}
	return s.repo.Update(ctx, id, record)
}

func (s *recordService) Delete(ctx context.Context, id int) error {
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

func (s *recordService) TotalPrice(ctx context.Context, userIDStr string, serviceName string, from time.Time, to time.Time) (int, error) {
	if to.Before(from) {
		return 0, ErrInvalidDates
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
