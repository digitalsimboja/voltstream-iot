# ---- Stage 1: build ----
FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /simulator ./cmd/main.go

# ---- Stage 2: run ----
FROM alpine:3.19

COPY --from=builder /simulator /simulator

# Defaults can be overridden via docker-compose command or docker run args
ENTRYPOINT ["/simulator"]
CMD ["--machines", "500", "--interval", "2s", "--target", "http://ingestion-api:8080/telemetry"]
