package saga

import (
    "encoding/json"
    "time"
)

type SagaState string

const (
    StateStarted   SagaState = "STARTED"
    StateCompleted SagaState = "COMPLETED"
    StateFailed    SagaState = "FAILED"
)

type SagaStep struct {
    StepName     string          `json:"step_name"`
    Compensation bool            `json:"compensation"`
    Payload      json.RawMessage `json:"payload"`
}

type SagaMessage struct {
    ID        string          `json:"id"`
    State     SagaState       `json:"state"`
    Steps     []SagaStep      `json:"steps"`
    Timestamp time.Time       `json:"timestamp"`
    Error     string         `json:"error,omitempty"`
}