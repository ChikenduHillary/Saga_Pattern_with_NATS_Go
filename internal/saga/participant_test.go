package saga

import (
	"errors"
	"testing"
)

func TestParticipantProcessSaga(t *testing.T) {
	// Create a participant without a real NATS connection for unit testing.
	p := NewParticipant(nil, "test-participant", "test-stream")

	t.Run("executes correct handler for step", func(t *testing.T) {
		handlerCalled := false
		p.RegisterHandler("TestStep", func(step SagaStep) error {
			handlerCalled = true
			return nil
		})

		sagaMsg := SagaMessage{
			Steps: []SagaStep{{StepName: "TestStep"}},
		}

		p.processSaga(sagaMsg)

		if !handlerCalled {
			t.Error("expected handler for 'TestStep' to be called, but it was not")
		}
	})

	t.Run("stops processing on handler error", func(t *testing.T) {
		var firstHandlerCalled, secondHandlerCalled bool

		p.RegisterHandler("FirstStep", func(step SagaStep) error {
			firstHandlerCalled = true
			return errors.New("something went wrong")
		})
		p.RegisterHandler("SecondStep", func(step SagaStep) error {
			secondHandlerCalled = true
			return nil
		})

		sagaMsg := SagaMessage{
			Steps: []SagaStep{
				{StepName: "FirstStep"},
				{StepName: "SecondStep"},
			},
		}

		p.processSaga(sagaMsg)

		if !firstHandlerCalled {
			t.Error("expected handler for 'FirstStep' to be called")
		}
		if secondHandlerCalled {
			t.Error("did not expect handler for 'SecondStep' to be called after first step failed")
		}
	})
}