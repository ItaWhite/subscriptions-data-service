FROM golang:1.26.4-alpine3.24 as builder
RUN apk add --no-cache git
RUN go install github.com/swaggo/swag/cmd/swag@v1.16.6
WORKDIR /code
ENTRYPOINT ["swag"]
