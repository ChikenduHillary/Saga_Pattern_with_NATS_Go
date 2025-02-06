package saga

import (
    "encoding/json"
    "fmt"
    "github.com/nats-io/nats.go"
    "log"
)

type ParticipantHandler func(step SagaStep) error

type Participant struct {
    js       nats.JetStreamContext
    name     string
    stream   string
    handlers map[string]ParticipantHandler
}

func NewParticipant(js nats.JetStreamContext, name, stream string) *Participant {
    return &Participant{
        js:       js,
        name:     name,
        stream:   stream,
        handlers: make(map[string]ParticipantHandler),
    }
}

func (p *Participant) RegisterHandler(stepName string, handler ParticipantHandler) {
    p.handlers[stepName] = handler
}

func (p *Participant) Start() error {
    _, err := p.js.Subscribe(
        fmt.Sprintf("%s.>", p.stream),
        func(msg *nats.Msg) {
            var sagaMsg SagaMessage
            if err := json.Unmarshal(msg.Data, &sagaMsg); err != nil {
                log.Printf("Error unmarshaling message: %v", err)
                return
            }

            for _, step := range sagaMsg.Steps {
                if handler, ok := p.handlers[step.StepName]; ok {
                    if err := handler(step); err != nil {
                        // Handle error and trigger compensation
                        p.handleStepError(sagaMsg, step, err)
                        return
                    }
                }
            }
        },
    )
    return err
}

func (p *Participant) handleStepError(saga SagaMessage, step SagaStep, err error) {
    // Implement compensation logic here
    log.Printf("Error in step %s: %v", step.StepName, err)
}