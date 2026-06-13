# VoltStream IoT — Industrial Battery Telemetry Platform

> Edge-to-cloud telemetry platform for industrial battery systems, built to mirror
> the exact architecture powering heavy-industry electrification at scale.

---

## Overview

VoltStream IoT is a high-throughput, edge-to-cloud telemetry system that ingests
real-time battery metrics from a simulated fleet of 500 off-road mining and
construction machines. It streams those metrics through an AWS data pipeline,
applies real-time anomaly detection, and surfaces predictive failure alerts
through a developer-facing dashboard.

The system is built as a full portfolio demonstration of the Go backend, AWS cloud,
Terraform IaC, and data engineering skills required for industrial-grade IoT platforms.

---

## Architecture

```
┌────────────────────────────────────────────────────────────┐
│             Simulated Battery Edge System                   │
│                  (Go Concurrent Engine)                     │
│       500 machines × telemetry packet every 2 seconds       │
└──────────────────────────┬─────────────────────────────────┘
                           │  MQTT / Secure Webhooks
                           ▼
                ┌─────────────────────────┐
                │  AWS API Gateway        │
                │  / AWS IoT Core         │
                └────────────┬────────────┘
                             │
                             ▼
                ┌─────────────────────────┐
                │  AWS Kinesis Data       │
                │  Streams / Amazon MSK   │
                └──────────┬──────────────┘
                           │
              ┌────────────┴─────────────┐
              ▼                          ▼
  ┌─────────────────────┐     ┌──────────────────────────┐
  │  AWS Lambda (Go)    │     │  AWS Glue Streaming      │
  │  Real-time Anomaly  │     │  / Amazon Timestream     │
  │  Detection          │     │  Long-term Analytics     │
  └──────────┬──────────┘     └─────────────┬────────────┘
             │                              │
             ▼                              ▼
  ┌─────────────────────┐     ┌──────────────────────────┐
  │  Slack / PagerDuty  │     │  Next.js Dashboard       │
  │  Alert Delivery     │     │  (WebSocket telemetry)   │
  └─────────────────────┘     └──────────────────────────┘
```

---

## Tech Stack

| Layer | Technology |
|---|---|
| Edge Simulator | Go (Golang) — concurrent goroutines |
| Ingestion API | Go — Gin or Fiber framework |
| Message Transport | MQTT, AWS IoT Core |
| Stream Ingestion | AWS Kinesis Data Streams / Amazon MSK (Kafka) |
| Real-time Processing | AWS Lambda (Go runtime) |
| Batch Analytics | AWS Glue, Amazon Timestream |
| Predictive AI | AWS Bedrock |
| Storage | Amazon Timestream, PostgreSQL (time-series) |
| Infrastructure as Code | Terraform — modular, least-privilege IAM |
| Frontend Dashboard | Next.js, React, Recharts, WebSockets |
| Developer CLI | Go — single-command environment provisioner |
| Observability | AWS CloudWatch, OpenTelemetry |
| CI/CD | GitHub Actions |

---

## Telemetry Data Payload

Each battery packet emitted by the Go edge simulator:

```json
{
  "battery_id": "BAT-0023-EXC",
  "machine_type": "Excavator",
  "voltage": 48.6,
  "temperature_celsius": 42.3,
  "state_of_charge_pct": 73.5,
  "current_amps": 120.4,
  "timestamp": "2026-06-13T10:42:00Z"
}
```

---

## Project Phases

### Phase 1 — Edge Simulator & Backend (Go)

- Concurrent Go CLI engine: 500 machines, one telemetry packet per machine every 2s
- Primary ingestion gateway: Go microservice using Gin or Fiber
- Concurrent upstream streaming using Go channels for backpressure safety

### Phase 2 — Infrastructure & Streaming Pipeline (AWS + Terraform)

- Modular Terraform for the full AWS stack: VPC, IAM, Kinesis, Lambda, Timestream
- Principle of least-privilege IAM roles; secrets in AWS Secrets Manager
- High-throughput stream: Kinesis Data Streams or Amazon MSK (Kafka)
- Lambda stream consumer (Go): drops corrupt packets, calculates moving averages,
  flags immediate anomalies (temperature > 80°C)

### Phase 3 — Storage & Predictive Analytics

- Amazon Timestream as the primary time-series store
- Async worker using AWS Bedrock to evaluate SoC degradation patterns
- Flags batteries whose SoC drops at an anomalous rate relative to operating temperature

### Phase 4 — Full-Stack Developer Interface

- Next.js dashboard: real-time fleet map, system health metrics, WebSocket telemetry charts
- Developer CLI: provisions a mocked cloud test environment matching production
  Terraform config with a single command — designed for Embedded/Hardware engineers

---

## Repository Structure

```
voltstream/
├── simulator/              # Go edge simulator engine
│   ├── cmd/main.go
│   └── internal/
│       ├── generator/      # telemetry packet generation
│       └── transport/      # MQTT / HTTP dispatch
├── ingestion-api/          # Go Gin ingestion gateway
│   ├── cmd/main.go
│   └── internal/
│       ├── handlers/
│       └── stream/
├── infra/                  # Terraform modules
│   ├── modules/
│   │   ├── vpc/
│   │   ├── iam/
│   │   ├── kinesis/
│   │   ├── lambda/
│   │   └── timestream/
│   └── environments/
│       ├── dev/
│       └── prod/
├── lambda/                 # Lambda stream processor (Go)
│   └── anomaly-detector/
│       └── main.go
├── analytics/              # Glue jobs + Bedrock predictive worker
├── dashboard/              # Next.js frontend
│   ├── components/
│   │   ├── FleetMap.tsx
│   │   ├── TelemetryChart.tsx
│   │   └── AlertPanel.tsx
│   └── pages/
├── cli/                    # Developer environment provisioner
│   └── cmd/voltctl/
├── docs/
│   └── architecture.md
└── README.md
```

---

## Getting Started

Prerequisites: Go 1.22+, Terraform 1.7+, Node.js 20+, AWS CLI configured

```bash
git clone https://github.com/digitalsimboja/voltstream-iot.git
cd voltstream-iot
```

---

## Makefile Commands

All common tasks are available via `make` from the repo root.

| Command | Description |
|---|---|
| `make sim` | Start the Go edge simulator — 500 machines, 2s interval |
| `make api` | Start the Go ingestion API on `:8080` |
| `make lambda-build` | Compile and zip the Lambda anomaly detector for deployment |
| `make test` | Run all Go tests (`go test ./...`) |
| `make lint` | Run `go vet` and Next.js ESLint |
| `make infra-dev-up` | Terraform init + apply for the dev environment |
| `make infra-dev-down` | Terraform destroy for the dev environment |
| `make dashboard` | Install dependencies and start the Next.js dashboard |
| `make voltctl-build` | Build the `voltctl` CLI binary to `bin/voltctl` |

### Quick start (local dev)

```bash
# 1. Start the ingestion API
make api

# 2. In a second terminal — start the simulator
make sim

# 3. In a third terminal — start the dashboard
make dashboard
```

### Infrastructure

```bash
# Provision the full dev environment on AWS
make infra-dev-up

# Tear it down
make infra-dev-down
```

### Build & test

```bash
# Run all Go unit tests
make test

# Lint Go code + dashboard
make lint

# Build Lambda artifact (Linux amd64 binary + zip)
make lambda-build

# Build the voltctl CLI
make voltctl-build
```

---

## Author

**Sunday Izuchukwu Mgbogu** — Technical Lead, Cloud & Backend Engineering  
[GitHub](https://github.com/digitalsimboja) | [LinkedIn](https://www.linkedin.com/in/sunday-mgbogu)
