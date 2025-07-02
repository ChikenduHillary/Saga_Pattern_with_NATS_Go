package saga

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/nats-io/nats.go"
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
		p.handleNatsMessage,
    )
    return err
}

func (p *Participant) handleNatsMessage(msg *nats.Msg) {
	var sagaMsg SagaMessage
	if err := json.Unmarshal(msg.Data, &sagaMsg); err != nil {
		log.Printf("Error unmarshaling saga message: %v", err)
		return
	}
	p.processSaga(sagaMsg)
}

// processSaga iterates through the steps of a saga and executes the registered handlers.
// This method is separate from the NATS message handling to allow for easier unit testing.
func (p *Participant) processSaga(sagaMsg SagaMessage) {
	for _, step := range sagaMsg.Steps {
		if handler, ok := p.handlers[step.StepName]; ok {
			if err := handler(step); err != nil {
				// Handle error and trigger compensation
				p.handleStepError(sagaMsg, step, err)
				return // Stop processing further steps on error
			}
		}
	}
}

func (p *Participant) handleStepError(saga SagaMessage, step SagaStep, err error) {
    // Implement compensation logic here
    log.Printf("Error in step %s: %v", step.StepName, err)
}