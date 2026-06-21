package model

import (
	"time"

	"github.com/google/uuid"
)

type Record struct {
	Id          int
	ServiceName string
	Price       int
	UserID      uuid.UUID
	StartDate   time.Time
	EndDate     *time.Time
}
