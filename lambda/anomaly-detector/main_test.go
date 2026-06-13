package main

import "testing"

func TestDetectAnomaly_HighTemperature(t *testing.T) {
	p := BatteryPacket{BatteryID: "BAT-0001-EXC", TemperatureCelsius: 85.0, StateOfChargePct: 50.0}
	if reason := detectAnomaly(p); reason == "" {
		t.Error("expected anomaly for temperature 85°C, got none")
	}
}

func TestDetectAnomaly_BoundaryTemperature(t *testing.T) {
	p := BatteryPacket{BatteryID: "BAT-0002-DRL", TemperatureCelsius: 80.0, StateOfChargePct: 50.0}
	if reason := detectAnomaly(p); reason != "" {
		t.Errorf("80°C is at threshold, not over — expected no anomaly, got: %s", reason)
	}
}

func TestDetectAnomaly_LowSoC(t *testing.T) {
	p := BatteryPacket{BatteryID: "BAT-0003-LDR", TemperatureCelsius: 40.0, StateOfChargePct: 5.0}
	if reason := detectAnomaly(p); reason == "" {
		t.Error("expected anomaly for SoC 5%, got none")
	}
}

func TestDetectAnomaly_Normal(t *testing.T) {
	p := BatteryPacket{BatteryID: "BAT-0004-HLR", TemperatureCelsius: 45.0, StateOfChargePct: 65.0}
	if reason := detectAnomaly(p); reason != "" {
		t.Errorf("expected no anomaly for normal packet, got: %s", reason)
	}
}
