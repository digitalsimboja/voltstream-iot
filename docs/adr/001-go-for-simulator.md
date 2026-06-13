# ADR 001 — Go for the Edge Simulator

**Status:** Accepted  
**Date:** 2026-06-13

## Context

The edge simulator must spawn 500 concurrent machine goroutines, each emitting
a telemetry packet every 2 seconds (~250 packets/sec total). The language choice
directly signals adaptability to Scania's preferred backend stack.

## Decision

Use **Go** for both the simulator and the ingestion API.

## Rationale

- Native goroutine concurrency maps directly to the "500 machines" model
  without the overhead of thread pools or async frameworks.
- Static typing catches schema drift at compile time, which matters when
  `BatteryPacket` must be consistent across the simulator, API, and Lambda.
- Scania's job description explicitly names Go — using it here demonstrates
  immediate adaptability rather than a hypothetical claim.
- The Gin framework provides production-grade HTTP routing with minimal
  ceremony, keeping the ingestion API readable.

## Alternatives Considered

- **Python (FastAPI):** Familiar but requires asyncio complexity for 500
  concurrent emitters; less idiomatic for CPU-bound tight loops.
- **Node.js:** Event loop is single-threaded; worker threads add complexity
  without benefiting the use case.

## Consequences

- Requires Go 1.22+ on the developer machine.
- Lambda must use the `provided.al2023` custom runtime, not a managed one.
