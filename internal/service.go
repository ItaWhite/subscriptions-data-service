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
	return s.repo.Create(ctx, record)
}

func (s *recordService) Update(ctx context.Context, id int, record Record) error {
	return s.repo.Update(ctx, id, record)
}

func (s *recordService) Delete(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}

func (s *recordService) TotalPrice(ctx context.Context, userIDStr string, serviceName string, from time.Time, to time.Time) (int, error) {
	return s.repo.GetTotalPrice(ctx, userIDStr, serviceName, from, to)
}
