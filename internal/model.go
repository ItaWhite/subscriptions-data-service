package internal

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type Record struct {
	Id          int         `json:"id"`
	ServiceName string      `json:"service_name"`
	Price       int         `json:"price"`
	UserID      pgtype.UUID `json:"user_id"`
	StartDate   pgtype.Date `json:"start_date"`
	EndDate     pgtype.Date `json:"end_date"`
}
