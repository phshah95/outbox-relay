package main

import (
	"context"
	"fmt"
	"outbox-relay/outbox"
	"outbox-relay/relay"
	"runtime"
	"time"
)

func main() {
	// Create a context with a cancel function.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Connect to the database.
	// db, err := sql.Open("postgres", "postgres://user:password@localhost/dbname")
	// if err != nil {
	// 	log.Fatalf("failed to connect to the database: %v", err)
	// }
	// defer db.Close()

	// Initialize the outbox component.
	// Initialize the relay component.
	relay := relay.NewRelay(ctx, time.Second*10)

	dbEvents := relay.DBEvents()
	relayedEvents := relay.RelayedEvents()

	// Start the relay component.

	go relay.Run()

	go func() {
		for {
			numGoRoutines := runtime.NumGoroutine()
			fmt.Println("Number of goroutines:", numGoRoutines)
			time.Sleep(time.Second * 3)
		}
	}()

	go func() {
		id := 0
		for {
			dbEvents <- outbox.OutboxEvent{
				Message:     []byte("yolo"),
				Id:          fmt.Sprintf("%d", id),
				Destination: "http://localhost:8080/relay",
			}
			id++
			time.Sleep(time.Millisecond * 300)
		}
	}()

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case event := <-relayedEvents:
				fmt.Println("Relayed event:", event)
			}
		}
	}()

	// outbox := outbox.NewOutbox(ctx, db, dbEvents, time.Second*10)
	// go func() {
	// 	if err := outbox.Run(); err != nil {
	// 		log.Fatalf("failed to start outbox: %v", err)
	// 	}
	// }()

	fmt.Println("Press ctrl + c to exit.")
	<-ctx.Done()
	fmt.Println("Exiting...")
}
