# ---- Stage 1: build ----
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Cache dependency downloads separately from source changes
COPY go.mod ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /api ./cmd/main.go

# ---- Stage 2: run ----
FROM alpine:3.19

COPY --from=builder /api /api

EXPOSE 8080
ENTRYPOINT ["/api"]
