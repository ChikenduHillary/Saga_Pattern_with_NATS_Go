package saga

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/nats-io/nats.go"
)

// TestMain runs setup before all tests in the package and teardown after.
// It ensures we have a running NATS server for integration tests.
func TestMain(m *testing.M) {
	// Check for NATS connection before running tests
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatalf("Cannot run integration tests: NATS server not available at %s. Please start a server with JetStream enabled.", nats.DefaultURL)
	}
	nc.Close()

	// Run all tests in the package
	os.Exit(m.Run())
}

func TestSagaEndToEnd(t *testing.T) {
	// --- Setup ---
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		t.Fatalf("Failed to connect to NATS: %v", err)
	}
	defer nc.Close()

	js, err := nc.JetStream()
	if err != nil {
		t.Fatalf("Failed to get JetStream context: %v", err)
	}

	// Use a unique stream name for each test run to ensure isolation
	streamName := fmt.Sprintf("order-saga-test-%d", time.Now().UnixNano())

	// Clean up the stream after the test
	defer func() {
		if err := js.DeleteStream(streamName); err != nil {
			t.Logf("Warning: failed to delete stream %s: %v", streamName, err)
		}
	}()

	// --- Create Coordinator and Participant ---
	coordinator, err := NewCoordinator(js, streamName)
	if err != nil {
		t.Fatalf("Failed to create coordinator: %v", err)
	}

	participant := NewParticipant(js, "test-order-service", streamName)

	// --- Test Logic ---
	// Create a channel to signal that the handler was successfully invoked
	handlerCh := make(chan SagaStep, 1)

	// Register a test handler that sends the received step to our channel
	participant.RegisterHandler("CreateOrder", func(step SagaStep) error {
		handlerCh <- step
		return nil
	})

	if err := participant.Start(); err != nil {
		t.Fatalf("Failed to start participant: %v", err)
	}

	// --- Action ---
	// Define and start the saga
	steps := []SagaStep{{StepName: "CreateOrder", Payload: json.RawMessage(`{"order_id": "test-123"}`)}}
	if _, err = coordinator.StartSaga(steps); err != nil {
		t.Fatalf("Failed to start saga: %v", err)
	}

	// --- Assert ---
	// Wait for the handler to be called, with a timeout
	select {
	case <-handlerCh:
		// Success! The message was received and processed.
	case <-time.After(5 * time.Second):
		t.Fatal("timed out waiting for participant handler to be called")
	}
}