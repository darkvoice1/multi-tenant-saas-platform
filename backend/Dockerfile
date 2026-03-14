FROM golang:1.25-alpine AS builder

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/saas-api ./cmd/api

FROM alpine:3.20

RUN apk add --no-cache ca-certificates
WORKDIR /app
RUN adduser -D -g '' appuser && mkdir -p /app/storage

COPY --from=builder /out/saas-api /app/saas-api
COPY migrations /app/migrations

ENV SERVER_ADDR=:8080
EXPOSE 8080

USER appuser
ENTRYPOINT ["/app/saas-api"]