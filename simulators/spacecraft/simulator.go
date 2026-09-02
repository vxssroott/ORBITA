package spacecraft

import (
	"context"
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/vxssroott/ORBITA/pkg/protocol"
)

type Config struct {
	SpacecraftID string
	Platform     string
	Interval     time.Duration
}

type Simulator struct {
	mu       sync.RWMutex
	config   Config
	sequence uint64
	started  time.Time
}

func New(cfg Config) (*Simulator, error) {
	if cfg.SpacecraftID == "" {
		return nil, fmt.Errorf("spacecraft id is required")
	}
	if cfg.Platform == "" {
		cfg.Platform = "ORBITA-SIM"
	}
	if cfg.Interval <= 0 {
		cfg.Interval = time.Second
	}

	return &Simulator{
		config:  cfg,
		started: time.Now().UTC(),
	}, nil
}

func (s *Simulator) NextTelemetry() protocol.TelemetryEnvelope {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.sequence++

	elapsed := time.Since(s.started).Seconds()

	return protocol.TelemetryEnvelope{
		SchemaVersion: "1.0",
		SpacecraftID:  s.config.SpacecraftID,
		Timestamp:     time.Now().UTC(),
		Sequence:      s.sequence,
		Parameters: map[string]any{
			"battery_voltage": 28.0 + math.Sin(elapsed/15.0)*0.4,
			"temperature":     24.0 + math.Sin(elapsed/10.0)*2.0,
			"signal_strength": 92.0 + math.Sin(elapsed/20.0)*3.0,
			"mode":            "nominal",
			"platform":        s.config.Platform,
		},
	}
}

func (s *Simulator) Run(ctx context.Context, emit func(protocol.TelemetryEnvelope) error) error {
	ticker := time.NewTicker(s.config.Interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			if err := emit(s.NextTelemetry()); err != nil {
				return err
			}
		}
	}
}
