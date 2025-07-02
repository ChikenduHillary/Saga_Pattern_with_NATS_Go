package order

import (
	"encoding/json"
	"saga-pattern/internal/saga"
	"testing"
)

func TestHandleCreateOrder(t *testing.T) {
	// Create a dummy participant for the unit test. It doesn't need a real
	// NATS connection, but it must be non-nil to allow handler registration.
	dummyParticipant := saga.NewParticipant(nil, "test-order-service", "test-saga")
	svc := NewOrderService(dummyParticipant)

	t.Run("valid payload", func(t *testing.T) {
		payload := CreateOrderPayload{
			OrderID: "order-456",
			Amount:  250.50,
			UserID:  "user-789",
		}
		payloadBytes, _ := json.Marshal(payload)

		step := saga.SagaStep{
			StepName: "CreateOrder",
			Payload:  payloadBytes,
		}

		err := svc.handleCreateOrder(step)
		if err != nil {
			t.Errorf("handleCreateOrder() returned an unexpected error: %v", err)
		}
	})

	t.Run("invalid payload", func(t *testing.T) {
		step := saga.SagaStep{
			StepName: "CreateOrder",
			Payload:  json.RawMessage(`{"invalid_json":}`),
		}

		err := svc.handleCreateOrder(step)
		if err == nil {
			t.Errorf("handleCreateOrder() was expected to return an error for invalid JSON, but it did not")
		}
	})
}