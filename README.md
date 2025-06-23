# Saga Pattern with NATS Streams in Go

Welcome! This project is a simple demonstration of the **Saga pattern** for distributed transactions in Go, using **NATS Streams (JetStream)** for communication between services. It’s designed to help beginners understand how to coordinate multi-step operations across microservices, ensuring data consistency even when things go wrong.

---

## What is a Saga?

A **Saga** is a design pattern used to manage long-running, multi-step business transactions in distributed systems. Instead of using a single database transaction, a saga breaks the process into a series of smaller, independent steps. Each step is handled by a different service and has a corresponding "compensating action" to undo its work if something goes wrong later in the process.

**Key Points:**

- Each step in a saga is a local transaction performed by a microservice.
- If a step fails, previously completed steps are compensated (undone) by running their compensating actions.
- Sagas help maintain data consistency across multiple services without distributed transactions.

**Example:**
Booking a trip might involve reserving a flight, a hotel, and a car. If the hotel reservation fails, the system cancels the flight and car reservations to keep everything consistent.

---

## Other Key Definitions

- **Distributed Transaction:** A transaction that involves multiple independent services or databases. Traditional database transactions (ACID) are hard to achieve in such systems, so patterns like Saga are used.
- **Compensating Action:** An operation that undoes the effects of a previously completed step in a saga. For example, if an order is created but payment fails, the compensating action would cancel the order.
- **Coordinator:** The component that manages the overall saga, tracks progress, and triggers compensating actions if needed.
- **Participant:** A service that performs a step in the saga and can also perform a compensating action if required.
- **NATS Streams (JetStream):** A messaging system used for communication between services. It allows services to publish and subscribe to events reliably.

---

## Table of Contents

- What is a Saga?
- Other Key Definitions
- Features
- How It Works
- Prerequisites
- Installation
- Running the Example
- Example Saga Flow
- Project Structure
- Example Code
- Troubleshooting
- Further Reading

---

## Features

- Demonstrates the Saga pattern for distributed transactions
- Uses [NATS Streams (JetStream)](https://docs.nats.io/nats-concepts/jetstream) for event-driven communication
- Includes a simple Order service as a Saga participant
- Handles compensation (undo) logic for failed steps
- Beginner-friendly code and explanations

---

## How It Works

1. The **Saga Coordinator** starts a saga (a multi-step transaction).
2. Each **Saga Participant** (e.g., Order service) performs its step and reports success or failure.
3. If any step fails, the coordinator triggers compensating actions to undo previous steps.
4. All communication between services happens via NATS Streams (JetStream).

**Analogy:**  
Imagine booking a trip: you reserve a flight, a hotel, and a car. If the hotel is unavailable, you cancel the flight and car. The Saga pattern helps automate this process in software.

---

## Prerequisites

- **Go**: Version 1.21 or higher. [Download Go](https://golang.org/dl/)
- **NATS Server**: [Download and install NATS Server](https://nats.io/download/nats-io/nats-server/)
- **Go Modules**:
  - `github.com/nats-io/nats.go`
  - `github.com/google/uuid`

---

## Installation

1. **Clone this repository**

   ```bash
   git clone https://github.com/ChikenduHillary/Saga_Pattern_with_NATS_Go.git
   cd Saga_Pattern_with_NATS_Go
   ```

2. **Install dependencies**
   ```bash
   go mod tidy
   ```

---

## Running the Example

1. **Start the NATS Server**  
   In a terminal, run:

   ```bash
   nats-server -js
   ```

   This starts NATS with JetStream enabled.

2. **Run the main program**

   ```bash
   go run ./cmd/main/main.go
   ```

3. **Observe the output**  
   The program will simulate a saga transaction and print logs for each step.

---

## Example Saga Flow

- The main program starts a saga to create an order.
- If any step fails, previous steps are compensated (undone).
- You can extend this by adding more services (e.g., payment, inventory).

---

## Project Structure

```
Saga_Pattern_with_NATS_Go/
├── cmd/
│   └── main/
│       └── main.go         # Entry point, starts the saga
└── internal/
    ├── order/
    │   └── service.go      # Order service logic
    └── saga/
        ├── coordinator.go  # Saga coordinator logic
        ├── participant.go  # Saga participant logic
        └── types.go        # Shared types
```

- **main.go**: Starts the saga and sets up NATS.
- **internal/saga/**: Contains the core logic for coordinating and participating in sagas.
- **internal/order/service.go**: Implements the order service as a saga participant.

---

## Example Code

### main.go

```go
package main

import (
    "encoding/json"
    "log"
    "github.com/nats-io/nats.go"
    "saga-nats-example/internal/saga"
    "saga-nats-example/internal/order"
)

func main() {
    // Connect to NATS
    nc, err := nats.Connect(nats.DefaultURL)
    if err != nil {
        log.Fatal(err)
    }
    defer nc.Close()

    // Create JetStream context
    js, err := nc.JetStream()
    if err != nil {
        log.Fatal(err)
    }

    // Set up the Saga Coordinator
    coordinator, err := saga.NewCoordinator(js, "order-saga")
    if err != nil {
        log.Fatal(err)
    }

    // Set up the Order Participant
    orderParticipant := saga.NewParticipant(js, "order-service", "order-saga")
    orderService := order.NewOrderService(orderParticipant)
    _ = orderService // Prevent unused variable warning

    // Start listening for saga messages
    if err := orderParticipant.Start(); err != nil {
        log.Fatal(err)
    }

    // Start a saga with a single step (CreateOrder)
    steps := []saga.SagaStep{
        {
            StepName: "CreateOrder",
            Payload: json.RawMessage(`{"order_id": "123"}`),
        },
    }
    if err := coordinator.StartSaga(steps); err != nil {
        log.Fatal(err)
    }
}
```

---

## Troubleshooting

- **NATS Server not running**: Make sure you started NATS with `nats-server -js`.
- **Port in use**: If port 4222 is busy, change it in your NATS config or use a different port.
- **Dependency issues**: Run `go mod tidy` to fix missing dependencies.

---

## Further Reading

- [Go Official Documentation](https://golang.org/doc/)
- [NATS Documentation](https://docs.nats.io/)
- [Saga Pattern (Microservices.io)](https://microservices.io/patterns/data/saga.html)

---

Happy coding! If you have questions, feel free to open an issue or ask for help.
