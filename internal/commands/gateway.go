package commands

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/vxssroott/ORBITA/pkg/protocol"
)

var (
	ErrInvalidCommand = errors.New("invalid command")
	ErrDuplicate      = errors.New("duplicate command")
)

type Transport interface {
	Send(context.Context, protocol.Command) error
}

type Gateway struct {
	mu        sync.Mutex
	transport Transport
	seen      map[string]struct{}
}

func NewGateway(transport Transport) *Gateway {
	return &Gateway{
		transport: transport,
		seen:      make(map[string]struct{}),
	}
}

func (g *Gateway) Submit(
	ctx context.Context,
	command protocol.Command,
) (protocol.CommandResult, error) {
	if command.ID == "" || command.SpacecraftID == "" || command.Type == "" {
		return protocol.CommandResult{}, ErrInvalidCommand
	}

	if command.CreatedAt.IsZero() {
		command.CreatedAt = time.Now().UTC()
	}

	g.mu.Lock()

	if _, exists := g.seen[command.ID]; exists {
		g.mu.Unlock()

		return protocol.CommandResult{
			CommandID: command.ID,
			Accepted:  false,
			Timestamp: time.Now().UTC(),
			Message:   ErrDuplicate.Error(),
		}, ErrDuplicate
	}

	g.seen[command.ID] = struct{}{}
	g.mu.Unlock()

	if g.transport == nil {
		return protocol.CommandResult{
			CommandID: command.ID,
			Accepted:  false,
			Timestamp: time.Now().UTC(),
			Message:   "command transport is unavailable",
		}, errors.New("command transport is unavailable")
	}

	if err := g.transport.Send(ctx, command); err != nil {
		return protocol.CommandResult{
			CommandID: command.ID,
			Accepted:  false,
			Timestamp: time.Now().UTC(),
			Message:   fmt.Sprintf("command rejected: %v", err),
		}, err
	}

	return protocol.CommandResult{
		CommandID:    command.ID,
		Accepted:     true,
		Acknowledged: true,
		Timestamp:    time.Now().UTC(),
		Message:      "command accepted",
	}, nil
}
