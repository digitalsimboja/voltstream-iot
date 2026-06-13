.PHONY: up down logs sim api lambda-build test lint infra-dev-up infra-dev-down dashboard voltctl-build

# --- docker compose ---
up:
	docker compose up --build

down:
	docker compose down

logs:
	docker compose logs -f

# --- simulator ---
sim:
	go run ./simulator/cmd/main.go --machines 500 --interval 2s

# --- ingestion API ---
api:
	go run ./ingestion-api/cmd/main.go

# --- lambda: build & package ---
lambda-build:
	GOOS=linux GOARCH=amd64 go build -o lambda/anomaly-detector/bootstrap ./lambda/anomaly-detector
	cd lambda/anomaly-detector && zip anomaly-detector.zip bootstrap

# --- tests ---
test:
	go test ./...

# --- linting ---
lint:
	go vet ./...
	cd dashboard && npm run lint

# --- terraform dev environment ---
infra-dev-up:
	terraform -chdir=infra/environments/dev init
	terraform -chdir=infra/environments/dev apply -auto-approve

infra-dev-down:
	terraform -chdir=infra/environments/dev destroy -auto-approve

# --- dashboard ---
dashboard:
	cd dashboard && npm install && npm run dev

# --- voltctl CLI ---
voltctl-build:
	go build -o bin/voltctl ./cli/cmd/voltctl
