package order

import (
    "encoding/json"
    "saga-nats-example/internal/saga"
)

type OrderService struct {
    sagaParticipant *saga.Participant
}

type CreateOrderPayload struct {
    OrderID  string  `json:"order_id"`
    Amount   float64 `json:"amount"`
    UserID   string  `json:"user_id"`
}

func NewOrderService(participant *saga.Participant) *OrderService {
    svc := &OrderService{
        sagaParticipant: participant,
    }

    participant.RegisterHandler("CreateOrder", svc.handleCreateOrder)
    participant.RegisterHandler("CancelOrder", svc.handleCancelOrder)

    return svc
}

func (s *OrderService) handleCreateOrder(step saga.SagaStep) error {
    var payload CreateOrderPayload
    if err := json.Unmarshal(step.Payload, &payload); err != nil {
        return err
    }
    
    // Implement order creation logic here
    return nil
}

func (s *OrderService) handleCancelOrder(step saga.SagaStep) error {
    var payload CreateOrderPayload
    if err := json.Unmarshal(step.Payload, &payload); err != nil {
        return err
    }
    
    // Implement order cancellation logic here
    return nil
}