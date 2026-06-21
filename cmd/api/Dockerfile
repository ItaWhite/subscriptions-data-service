FROM golang:1.26.4-alpine3.24 AS builder
WORKDIR /app
COPY go.* ./
RUN go mod download
COPY ../.. .
RUN go build -o /app/cmd/api/main /app/cmd/api

FROM alpine:3.24
WORKDIR /app
COPY --from=builder /app/cmd/api/main /app
EXPOSE 8080
CMD ["./main"]