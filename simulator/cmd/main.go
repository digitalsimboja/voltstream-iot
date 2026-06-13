package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/digitalsimboja/voltstream-iot/simulator/internal/generator"
	"github.com/digitalsimboja/voltstream-iot/simulator/internal/transport"
)

func main() {
	machines := flag.Int("machines", 10, "number of machines to simulate")
	interval := flag.Duration("interval", 2*time.Second, "telemetry emit interval per machine")
	target := flag.String("target", "http://localhost:8080/telemetry", "ingestion API endpoint")
	flag.Parse()

	slog.Info("starting VoltStream simulator",
		"machines", *machines,
		"interval", *interval,
		"target", *target,
	)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	client := transport.NewHTTPClient(*target)

	for i := 0; i < *machines; i++ {
		go runMachine(i, *interval, client)
	}

	<-stop
	fmt.Println("\nshutting down simulator")
}

func runMachine(id int, interval time.Duration, client *transport.HTTPClient) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		packet := generator.GeneratePacket(id)
		if err := client.Send(packet); err != nil {
			slog.Error("failed to send packet", "machine_id", id, "error", err)
		} else {
			slog.Info("packet sent",
				"battery_id", packet.BatteryID,
				"temp", packet.TemperatureCelsius,
				"soc", packet.StateOfChargePct,
			)
		}
	}
}
