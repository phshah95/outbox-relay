package relay

import (
	"context"
	"fmt"
	"time"

	"outbox-relay/outbox"
)

// Relay represents the relay component of the outbox relay.
type Relay struct {
	ctx           context.Context
	dbEvents      chan outbox.OutboxEvent
	relayedEvents chan outbox.OutboxEvent
	endpoint      string
	timeout       time.Duration
}

// NewRelay creates a new instance of the Relay component.
func NewRelay(ctx context.Context, timeout time.Duration) *Relay {
	return &Relay{
		ctx:           ctx,
		dbEvents:      make(chan outbox.OutboxEvent, 10),
		relayedEvents: make(chan outbox.OutboxEvent),
		endpoint:      "http://localhost:8080/relay",
		timeout:       timeout,
	}
}

// Start starts the relay component and listens for incoming messages on the message channel.
func (r *Relay) Run() error {
	for {
		select {
		case <-r.ctx.Done():
			return nil
		// read all messages from the channel and send them to the relay endpoint.
		case event := <-r.dbEvents:
			if err := r.relayEvent(event); err != nil {
				return fmt.Errorf("failed to relay event: %v", err)
			}
			r.relayedEvents <- event
		}
	}
}

func (r *Relay) DBEvents() chan outbox.OutboxEvent {
	return r.dbEvents
}

func (r *Relay) RelayedEvents() chan outbox.OutboxEvent {
	return r.relayedEvents
}

// relayEvent sends a message to the relay endpoint.
func (r *Relay) relayEvent(event outbox.OutboxEvent) error {
	// Create a context with timeout for the request.
	// ctx, cancel := context.WithTimeout(r.ctx, r.timeout)
	// defer cancel()
	fmt.Println("Relaying event:", event)
	time.Sleep(time.Second * 10)
	return nil
}
