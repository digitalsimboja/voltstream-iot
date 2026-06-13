package generator

import (
	"fmt"
	"math/rand"
	"time"
)

// MachineType identifies the category of industrial machine.
type MachineType string

const (
	Excavator    MachineType = "Excavator"
	Drill        MachineType = "Drill"
	Loader       MachineType = "Loader"
	Hauler       MachineType = "Hauler"
	MarineVessel MachineType = "MarineVessel"
)

var machineTypes = []MachineType{Excavator, Drill, Loader, Hauler, MarineVessel}

// BatteryPacket is the telemetry payload emitted by each simulated machine.
type BatteryPacket struct {
	BatteryID          string      `json:"battery_id"`
	MachineType        MachineType `json:"machine_type"`
	Voltage            float64     `json:"voltage"`
	TemperatureCelsius float64     `json:"temperature_celsius"`
	StateOfChargePct   float64     `json:"state_of_charge_pct"`
	CurrentAmps        float64     `json:"current_amps"`
	Timestamp          time.Time   `json:"timestamp"`
}

// GeneratePacket produces a realistic telemetry packet for the given machine index.
// Ranges are intentionally wide enough to occasionally breach anomaly thresholds.
func GeneratePacket(machineID int) BatteryPacket {
	mt := machineTypes[machineID%len(machineTypes)]
	return BatteryPacket{
		BatteryID:          fmt.Sprintf("BAT-%04d-%s", machineID, shortCode(mt)),
		MachineType:        mt,
		Voltage:            trunc(44.0 + rand.Float64()*8.0),   // 44–52 V
		TemperatureCelsius: trunc(25.0 + rand.Float64()*65.0),  // 25–90°C (spans anomaly threshold of 80)
		StateOfChargePct:   trunc(5.0 + rand.Float64()*90.0),   // 5–95 % (spans low-SoC threshold of 10)
		CurrentAmps:        trunc(80.0 + rand.Float64()*100.0), // 80–180 A
		Timestamp:          time.Now().UTC(),
	}
}

func shortCode(mt MachineType) string {
	codes := map[MachineType]string{
		Excavator:    "EXC",
		Drill:        "DRL",
		Loader:       "LDR",
		Hauler:       "HLR",
		MarineVessel: "MRN",
	}
	if c, ok := codes[mt]; ok {
		return c
	}
	return "UNK"
}

// trunc rounds to one decimal place.
func trunc(f float64) float64 {
	return float64(int(f*10)) / 10
}
