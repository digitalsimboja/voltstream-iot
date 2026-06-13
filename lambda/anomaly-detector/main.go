package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

const (
	maxTempCelsius = 80.0
	minSoCPct      = 10.0
)

// BatteryPacket mirrors the ingestion API schema.
type BatteryPacket struct {
	BatteryID          string  `json:"battery_id"`
	MachineType        string  `json:"machine_type"`
	Voltage            float64 `json:"voltage"`
	TemperatureCelsius float64 `json:"temperature_celsius"`
	StateOfChargePct   float64 `json:"state_of_charge_pct"`
	CurrentAmps        float64 `json:"current_amps"`
}

// Anomaly holds a flagged packet and the threshold it breached.
type Anomaly struct {
	Packet BatteryPacket `json:"packet"`
	Reason string        `json:"reason"`
}

func handler(ctx context.Context, event events.KinesisEvent) error {
	var anomalies []Anomaly

	for _, record := range event.Records {
		var packet BatteryPacket
		if err := json.Unmarshal(record.Kinesis.Data, &packet); err != nil {
			slog.Warn("dropping malformed record", "error", err, "shard", record.EventSourceArn)
			continue
		}

		if reason := detectAnomaly(packet); reason != "" {
			anomalies = append(anomalies, Anomaly{Packet: packet, Reason: reason})
		}
	}

	for _, a := range anomalies {
		slog.Warn("anomaly detected",
			"battery_id", a.Packet.BatteryID,
			"reason", a.Reason,
			"temp_c", a.Packet.TemperatureCelsius,
			"soc_pct", a.Packet.StateOfChargePct,
		)
		// TODO: write anomaly record to Amazon Timestream
		// TODO: publish alert to SNS topic → Slack / PagerDuty
	}

	slog.Info("batch processed", "total", len(event.Records), "anomalies", len(anomalies))
	return nil
}

// detectAnomaly returns a human-readable reason string if the packet breaches a threshold.
// Returns an empty string if the packet is within normal operating parameters.
func detectAnomaly(p BatteryPacket) string {
	if p.TemperatureCelsius > maxTempCelsius {
		return fmt.Sprintf("temperature %.1f°C exceeds %.0f°C threshold", p.TemperatureCelsius, maxTempCelsius)
	}
	if p.StateOfChargePct < minSoCPct {
		return fmt.Sprintf("SoC %.1f%% is below %.0f%% minimum", p.StateOfChargePct, minSoCPct)
	}
	// TODO: detect anomalous SoC drop rate using a per-battery sliding window
	return ""
}

func main() {
	lambda.Start(handler)
}
