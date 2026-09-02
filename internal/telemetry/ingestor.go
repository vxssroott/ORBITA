package telemetry

import (
	"context"
	"errors"

	"github.com/vxssroott/ORBITA/pkg/protocol"
)

type Handler interface {
	Handle(context.Context, protocol.TelemetryEnvelope) error
}

type Ingestor struct {
	validator *Validator
	store     *Store
	handler   Handler
}

func NewIngestor(validator *Validator, store *Store, handler Handler) *Ingestor {
	return &Ingestor{
		validator: validator,
		store:     store,
		handler:   handler,
	}
}

func (i *Ingestor) Ingest(ctx context.Context, envelope protocol.TelemetryEnvelope) error {
	if i.validator == nil {
		return errors.New("telemetry validator is not configured")
	}

	if err := i.validator.Validate(&envelope); err != nil {
		return err
	}

	if i.store != nil {
		i.store.Put(envelope)
	}

	if i.handler != nil {
		return i.handler.Handle(ctx, envelope)
	}

	return nil
}
