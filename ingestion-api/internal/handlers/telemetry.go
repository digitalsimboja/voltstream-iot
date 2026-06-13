package handlers

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// BatteryPacket mirrors the simulator payload schema.
type BatteryPacket struct {
	BatteryID          string    `json:"battery_id"          binding:"required"`
	MachineType        string    `json:"machine_type"        binding:"required"`
	Voltage            float64   `json:"voltage"`
	TemperatureCelsius float64   `json:"temperature_celsius"`
	StateOfChargePct   float64   `json:"state_of_charge_pct"`
	CurrentAmps        float64   `json:"current_amps"`
	Timestamp          time.Time `json:"timestamp"`
}

// IngestTelemetry validates and accepts a battery telemetry packet.
// POST /telemetry
func IngestTelemetry(c *gin.Context) {
	var packet BatteryPacket
	if err := c.ShouldBindJSON(&packet); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	slog.Info("packet received",
		"battery_id", packet.BatteryID,
		"machine_type", packet.MachineType,
		"temp_c", packet.TemperatureCelsius,
		"soc_pct", packet.StateOfChargePct,
	)

	// TODO: publish packet to Kinesis stream via stream.Publisher
	// if err := streamPublisher.Publish(c.Request.Context(), packet); err != nil {
	// 	slog.Error("kinesis publish failed", "error", err)
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "stream unavailable"})
	// 	return
	// }

	// TODO: broadcast packet to connected WebSocket clients

	c.JSON(http.StatusAccepted, gin.H{
		"status":     "accepted",
		"battery_id": packet.BatteryID,
	})
}
