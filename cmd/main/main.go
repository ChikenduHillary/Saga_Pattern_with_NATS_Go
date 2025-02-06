package main

import (
    "encoding/json"
    "log"
    "github.com/nats-io/nats.go"
    "saga-nats-example/internal/saga"
    "saga-nats-example/internal/order"
)

func main() {
    nc, err := nats.Connect(nats.DefaultURL)
    if err != nil {
        log.Fatal(err)
    }
    defer nc.Close()

    js, err := nc.JetStream()
    if err != nil {
        log.Fatal(err)
    }

    coordinator, err := saga.NewCoordinator(js, "order-saga")
    if err != nil {
        log.Fatal(err)
    }

    orderParticipant := saga.NewParticipant(js, "order-service", "order-saga")
    orderService := order.NewOrderService(orderParticipant)

    if err := orderParticipant.Start(); err != nil {
        log.Fatal(err)
    }

    // Example: Start a saga
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

    sagaMsg, err := coordinator.StartSaga(steps)
    if err != nil {
        log.Fatal(err)
    }

    log.Printf("Started saga with ID: %s", sagaMsg.ID)

    select {}
}