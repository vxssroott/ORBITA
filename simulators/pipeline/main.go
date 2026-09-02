package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/vxssroott/ORBITA/internal/events"
	"github.com/vxssroott/ORBITA/internal/rules"
	"github.com/vxssroott/ORBITA/internal/state"
	"github.com/vxssroott/ORBITA/internal/telemetry"
	"github.com/vxssroott/ORBITA/pkg/protocol"
	"github.com/vxssroott/ORBITA/simulators/spacecraft"
)

func floatPtr(value float64) *float64 {
	return &value
}

func main() {
	fmt.Println("============================================================")
	fmt.Println("ORBITA - SPACECRAFT TELEMETRY PIPELINE")
	fmt.Println("============================================================")
	fmt.Println("")

	ctx, cancel := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer cancel()

	sc, err := spacecraft.New(spacecraft.Config{
		SpacecraftID: "NIGCOMSAT-SIM-01",
		Platform:     "ORBITA-SIMULATED-SPACECRAFT",
		Interval:     300 * time.Millisecond,
	})
	if err != nil {
		panic(err)
	}

	validator := telemetry.NewValidator()
	stateEngine := state.NewEngine()
	eventEngine := events.NewEngine()

	ruleEngine := rules.NewEngine([]rules.Rule{
		{
			Name:      "temperature-critical-high",
			Parameter: "temperature",
			Max:       floatPtr(80),
			Severity:  rules.SeverityCritical,
			EventType: "thermal_anomaly",
			Message:   "spacecraft temperature exceeds safe operating threshold",
		},
		{
			Name:      "battery-critical-low",
			Parameter: "battery_voltage",
			Min:       floatPtr(20),
			Severity:  rules.SeverityCritical,
			EventType: "battery_anomaly",
			Message:   "spacecraft battery voltage is below safe operating threshold",
		},
		{
			Name:      "signal-degraded-low",
			Parameter: "signal_strength",
			Min:       floatPtr(20),
			Severity:  rules.SeverityWarning,
			EventType: "communications_degraded",
			Message:   "spacecraft signal strength is below operational threshold",
		},
	})

	fmt.Println("[BOOT] Spacecraft simulator: ONLINE")
	fmt.Println("[BOOT] Telemetry validator: ONLINE")
	fmt.Println("[BOOT] State engine: ONLINE")
	fmt.Println("[BOOT] Rule engine: ONLINE")
	fmt.Println("[BOOT] Event engine: ONLINE")
	fmt.Println("")
	fmt.Println("STREAM STARTED")
	fmt.Println("")

	count := 0
	anomalyInjected := false

	err = sc.Run(ctx, func(envelope protocol.TelemetryEnvelope) error {
		count++

		if count == 6 {
			anomalyInjected = true

			fmt.Println("")
			fmt.Println("!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!")
			fmt.Println("ANOMALY INJECTION")
			fmt.Println("!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!")

			envelope.Parameters["temperature"] = 95.0
			envelope.Parameters["battery_voltage"] = 17.5
			envelope.Parameters["signal_strength"] = 12.0
			envelope.Parameters["mode"] = "safe-mode"

			fmt.Println("[INJECT] temperature=95.00C")
			fmt.Println("[INJECT] battery_voltage=17.50V")
			fmt.Println("[INJECT] signal_strength=12.00%")
			fmt.Println("[INJECT] mode=safe-mode")
			fmt.Println("")
		}

		if err := validator.Validate(&envelope); err != nil {
			return fmt.Errorf("telemetry validation: %w", err)
		}

		stateValue, err := stateEngine.Apply(envelope)
		if err != nil {
			return fmt.Errorf("state engine: %w", err)
		}

		fmt.Printf(
			"[TELEMETRY] seq=%03d spacecraft=%s temp=%.2fC battery=%.2fV signal=%.2f%% health=%s\n",
			envelope.Sequence,
			envelope.SpacecraftID,
			envelope.Parameters["temperature"],
			envelope.Parameters["battery_voltage"],
			envelope.Parameters["signal_strength"],
			stateValue.Health,
		)

		operationalEvents := ruleEngine.Evaluate(
			envelope.SpacecraftID,
			envelope.Parameters,
		)

		for _, event := range operationalEvents {
			emitted := eventEngine.Emit(event)

			fmt.Printf(
				"[EVENT] id=%s type=%s severity=%s spacecraft=%s message=%s\n",
				emitted.ID,
				emitted.Type,
				emitted.Severity,
				emitted.SpacecraftID,
				emitted.Message,
			)
		}

		if count >= 10 {
			fmt.Println("")
			fmt.Println("============================================================")

			if anomalyInjected && len(eventEngine.Recent(100)) >= 3 {
				fmt.Println("ANOMALY DETECTION VERIFIED")
			} else {
				fmt.Println("ANOMALY DETECTION FAILED")
			}

			fmt.Println("============================================================")
			fmt.Println("")
			fmt.Printf("[RESULT] Telemetry packets processed: %d\n", count)
			fmt.Printf("[RESULT] Operational events captured: %d\n", len(eventEngine.Recent(100)))
			fmt.Printf("[RESULT] Final spacecraft health: %s\n", stateValue.Health)
			fmt.Println("")

			cancel()
		}

		return nil
	})

	if err != nil && err != context.Canceled {
		fmt.Printf("[ERROR] pipeline stopped: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("PIPELINE SHUTDOWN: CLEAN")
}
