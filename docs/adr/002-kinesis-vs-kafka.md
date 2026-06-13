# ADR 002 — AWS Kinesis Data Streams vs Amazon MSK (Kafka)

**Status:** Accepted (Kinesis chosen for initial build; MSK path documented)  
**Date:** 2026-06-13

## Context

The streaming layer must ingest ~250 packets/sec from the ingestion API and
fan out to Lambda (real-time) and Glue (batch analytics). Both Kinesis and
MSK (managed Kafka) are viable.

## Decision

Use **AWS Kinesis Data Streams** for the portfolio build.

## Rationale

| Criteria | Kinesis | MSK (Kafka) |
|---|---|---|
| Setup complexity | Low (Terraform in ~20 lines) | High (VPC, brokers, ZooKeeper/KRaft config) |
| Lambda integration | Native event source mapping | Requires MSK trigger + VPC peering |
| Ops overhead | Fully managed | Broker patching, scaling, storage management |
| Cost at low throughput | ~$0.015/shard-hour | Broker compute always running |
| Ecosystem fit | AWS-native | Portable; better for multi-cloud |
| Throughput ceiling | 1 MB/s per shard | Effectively unlimited with partition scaling |

At 250 packets/sec (~25 KB/s assuming 100-byte payloads), a single Kinesis shard
is sufficient. The project is AWS-native, so the managed integration wins.

## Migration Path to MSK

If requirements grow beyond Kinesis limits or multi-cloud portability becomes
necessary, the `ingestion-api/internal/stream` package is abstracted behind a
`Publisher` interface. Swapping the implementation to a Kafka producer
(`segmentio/kafka-go`) requires no changes to the handler layer.

## Consequences

- Kinesis record retention is capped at 7 days (vs Kafka's configurable
  unlimited). Acceptable for this use case.
- Records are base64-encoded in Lambda events — the anomaly detector decodes
  via `record.Kinesis.Data` (already handled in `main.go`).
