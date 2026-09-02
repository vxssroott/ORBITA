package telemetry

import (
	"time"

	"github.com/vxssroott/ORBITA/pkg/protocol"
)

type Generator struct {
	spacecraftID string
	sequence     uint64
}

func NewGenerator(spacecraftID string) *Generator {
	return &Generator{
		spacecraftID: spacecraftID,
	}
}

func (g *Generator) Next(parameters map[string]any) protocol.TelemetryEnvelope {
	g.sequence++

	return protocol.TelemetryEnvelope{
		SchemaVersion: "1.0",
		SpacecraftID:  g.spacecraftID,
		Timestamp:     time.Now().UTC(),
		Sequence:      g.sequence,
		Parameters:    parameters,
	}
}
