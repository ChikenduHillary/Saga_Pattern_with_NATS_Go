package main

// This is the entry point of the application.
// It demonstrates how to start a Saga transaction using NATS Streams (JetStream).

import (
	"encoding/json"
	"log"

	// Import the NATS Go client
	"github.com/nats-io/nats.go"
	// Import the saga logic and order service
	"saga-pattern/internal/order"
	"saga-pattern/internal/saga"
)

func main() {
    // Connect to the local NATS server (default URL is nats://localhost:4222)
    nc, err := nats.Connect(nats.DefaultURL)
    if err != nil {
        log.Fatal(err)
    }
    defer nc.Close()

    // Create a JetStream context for advanced NATS features
    js, err := nc.JetStream()
    if err != nil {
        log.Fatal(err)
    }

    // Create a Saga Coordinator to manage the saga workflow
    coordinator, err := saga.NewCoordinator(js, "order-saga")
    if err != nil {
        log.Fatal(err)
    }

    // Create a Saga Participant for the order service
    orderParticipant := saga.NewParticipant(js, "order-service", "order-saga")
    // Initialize the order service with the participant
    orderService := order.NewOrderService(orderParticipant)
    _ = orderService // Prevent unused variable warning (for beginners)

    // Start listening for saga messages in the order participant
    if err := orderParticipant.Start(); err != nil {
        log.Fatal(err)
    }

    // Example: Start a saga with a single step (CreateOrder)
    steps := []saga.SagaStep{
        {
            StepName: "CreateOrder",
            Payload: json.RawMessage(`{
                "order_id": "123",
                "amount": 100.0,
                "user_id": "user-123"
            }`),
        },
    }

    // Start the saga using the coordinator
    sagaMsg, err := coordinator.StartSaga(steps)
    if err != nil {
        log.Fatal(err)
    }

    log.Printf("Started saga with ID: %s", sagaMsg.ID)

    // Block the main goroutine forever (or until the program is terminated)
    select {}
}