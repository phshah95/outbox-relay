package outbox

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

// Outbox represents an instance of the outbox component.
type Outbox struct {
	ctx             context.Context
	db              *sql.DB
	dbMessages      chan OutboxEvent
	interval        time.Duration
	relayedMessages chan OutboxEvent
	stopped         bool
	stopSignal      chan struct{}
}

// NewOutbox returns a new instance of the outbox component.
func NewOutbox(ctx context.Context, db *sql.DB, outboxCh chan OutboxEvent, interval time.Duration, sentCh chan OutboxEvent) *Outbox {
	return &Outbox{
		ctx:             ctx,
		db:              db,
		dbMessages:      outboxCh,
		interval:        interval,
		relayedMessages: sentCh,
		stopped:         false,
		stopSignal:      make(chan struct{}),
	}
}

func (o *Outbox) Run() error {
	ticker := time.NewTicker(o.interval)
	defer ticker.Stop()

	for {
		select {
		case <-o.ctx.Done():
			return nil

		case <-ticker.C:
			rows, err := o.db.Query("SELECT message FROM outbox WHERE sent = false")
			if err != nil {
				return fmt.Errorf("failed to read messages from the outbox table: %v", err)
			}
			defer rows.Close()

			for rows.Next() {
				var message OutboxEvent
				if err := rows.Scan(&message); err != nil {
					return fmt.Errorf("failed to scan message from the outbox table: %v", err)
				}
				o.dbMessages <- message
			}

		case message := <-o.relayedMessages:
			if _, err := o.db.Exec("UPDATE outbox SET sent = true WHERE message = $1", message); err != nil {
				return fmt.Errorf("failed to mark message as sent in the database: %v", err)
			}
		}
	}
}
