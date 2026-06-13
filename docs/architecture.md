# VoltStream IoT — Architecture Reference

## System Overview

VoltStream IoT is a four-phase, edge-to-cloud telemetry system. Each phase is
independently deployable and represents a distinct architectural boundary.

---

## Data Flow

```
1. EDGE LAYER (Go simulator)
   500 goroutines × 1 packet per 2s → ~250 packets/sec peak

2. TRANSPORT
   HTTP POST → ingestion API  (swap to MQTT / AWS IoT Core for production hardware)

3. INGESTION API (Go / Gin)
   Validates schema → publishes to Kinesis via PutRecord

4. STREAM (Kinesis Data Streams)
   Ordered, partitioned by battery_id
   Retention: 24h (dev), 48h (prod)

5. PROCESSING (Lambda Go)
   Reads Kinesis records → filters corrupt data →
   detects threshold anomalies → writes to Timestream

6. STORAGE (Amazon Timestream)
   Two tables:
   - battery-metrics   → full telemetry, 24h memory / 90d magnetic
   - anomalies         → flagged events, 24h memory / 365d magnetic

7. ANALYTICS (AWS Glue + Bedrock)
   Glue: moving averages, cleaning, long-term aggregations
   Bedrock: per-battery SoC degradation classification

8. PRESENTATION (Next.js)
   WebSocket feed from ingestion API →
   AlertPanel + FleetMap + TelemetryChart
```

---

## Key Design Decisions

See [ADR 001](adr/001-go-for-simulator.md) — Go for the edge simulator  
See [ADR 002](adr/002-kinesis-vs-kafka.md) — Kinesis vs Kafka trade-off

---

## Scaling Considerations

| Bottleneck | Current Setting | Scale Path |
|---|---|---|
| Kinesis shards | 1 (dev), 4 (prod) | Switch to ON_DEMAND mode |
| Lambda concurrency | default | Set reserved concurrency + DLQ |
| Timestream writes | single-record | Batch PutRecords (up to 100/call) |
| Dashboard WebSocket | in-process | Move to Redis Pub/Sub fanout |

---

## Security Boundaries

- All Lambda and Glue roles follow least-privilege IAM
- Kinesis stream encrypted at rest with KMS (`alias/aws/kinesis`)
- Secrets (API keys, SNS ARNs) stored in AWS Secrets Manager — never in env vars
- VPC private subnets isolate Lambda and Glue from public internet
