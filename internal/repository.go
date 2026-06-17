package internal

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
)

type recordRepository struct {
	db *pgxpool.Pool
}

func NewRecordRepository(pool *pgxpool.Pool) *recordRepository {
	return &recordRepository{
		db: pool,
	}
}

func (r *recordRepository) InitSchema() error {
	_, err := r.db.Exec(context.Background(), `
create table if not exists records(
    id int generated always as identity primary key,
    service_name varchar(255) not null,
    price int not null,
    user_id uuid not null,
    start_date date not null,
    end_date date
);`)
	if err != nil {
		return err
	}

	_, err = r.db.Exec(context.Background(), `
insert into records (service_name, price, user_id, start_date, end_date) values ('Yandex Plus', 400, '60601fee-2bf1-4721-ae6f-7636e79a0cba', '2026-01-01', '2026-07-01');`)
	return err
}

func (r *recordRepository) DropSchema() error {
	_, err := r.db.Exec(context.Background(), `
drop table if exists records;
`)
	return err
}

func (r *recordRepository) GetAll() ([]Record, error) {
	rows, err := r.db.Query(context.Background(), "select * from records;")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var recordList []Record
	for rows.Next() {
		var record Record
		err = rows.Scan(&record.Id, &record.ServiceName, &record.Price, &record.UserID, &record.StartDate, &record.EndDate)
		if err != nil {
			return nil, err
		}
		recordList = append(recordList, record)
	}
	return recordList, nil
}

func (r *recordRepository) GetByID(id int) (Record, error) {
	var record Record
	err := r.db.QueryRow(context.Background(), "select * from records where id=$1", id).
		Scan(&record.Id, &record.ServiceName, &record.Price, &record.UserID, &record.StartDate, &record.EndDate)
	if err != nil {
		return Record{}, err
	}
	return record, err
}

func (r *recordRepository) Create(record Record) (Record, error) {
	err := r.db.QueryRow(context.Background(), `
insert into records (service_name, price, user_id, start_date, end_date) 
values ($1, $2, $3, $4, $5) returning id;`, record.ServiceName, record.Price, record.UserID, record.StartDate, record.EndDate).Scan(&record.Id)
	if err != nil {
		return Record{}, err
	}
	return record, nil
}

func (r *recordRepository) Update(id int, record Record) error {
	cmd, err := r.db.Exec(context.Background(), "update records set service_name=$1, price=$2, user_id=$3, start_date=$4, end_date=$5 where id=$6",
		record.ServiceName, record.Price, record.UserID, record.StartDate, record.EndDate, id)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return errors.New("record not found")
	}
	record.Id = id
	return err
}

func (r *recordRepository) Delete(id int) error {
	_, err := r.db.Exec(context.Background(), "delete from records where id=$1", id)
	if err != nil {
		return err
	}
	return nil
}
