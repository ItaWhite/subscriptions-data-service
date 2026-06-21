include .env
export

export PROJECT_ROOT=$(shell pwd)

env-up:
	docker compose up -d subscriptions-postgres

env-down:
	docker compose down subscriptions-postgres

env-cleanup:
	docker compose down subscriptions-postgres
	rm -rf out/pgdata

migrate-create:
	docker compose run --rm subscriptions-postgres-migrate create -ext sql -dir /migrations -seq create_records

migrate-up:
	docker compose run --rm subscriptions-postgres-migrate -path /migrations -database postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@subscriptions-postgres:5432/${POSTGRES_DB}?sslmode=disable up

migrate-down:
	docker compose run --rm subscriptions-postgres-migrate -path /migrations -database postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@subscriptions-postgres:5432/${POSTGRES_DB}?sslmode=disable down

service-deploy:
	docker compose up -d --build subscriptions-service

service-stop:
	docker compose down subscriptions-service

swagger-create:
	docker compose run --rm swagger init -g cmd/api/main.go -o docs --parseInternal --parseDependency
