package saga

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
)

type Coordinator struct {
    js         nats.JetStreamContext
    sagaStream string
}

func NewCoordinator(js nats.JetStreamContext, stream string) (*Coordinator, error) {
    _, err := js.AddStream(&nats.StreamConfig{
        Name:     stream,
        Subjects: []string{fmt.Sprintf("%s.>", stream)},
    })
    if err != nil && err != nats.ErrStreamNameAlreadyInUse {
        return nil, err
    }

    return &Coordinator{
        js:         js,
        sagaStream: stream,
    }, nil
}

func (c *Coordinator) StartSaga(steps []SagaStep) (*SagaMessage, error) {
    sagaID := uuid.New().String()
    msg := &SagaMessage{
        ID:        sagaID,
        State:     StateStarted,
        Steps:     steps,
        Timestamp: time.Now(),
    }

    data, err := json.Marshal(msg)
    if err != nil {
        return nil, err
    }

    _, err = c.js.Publish(
        fmt.Sprintf("%s.%s", c.sagaStream, sagaID),
        data,
    )
    return msg, err
}