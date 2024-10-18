package event_bus

import "fmt"

type EventBus struct {
	subscribers map[string][]chan<- string
}

func NewEventBus() *EventBus {
	return &EventBus{
		subscribers: make(map[string][]chan<- string),
	}
}

func (eventBus *EventBus) Publish(eventType string, msg string) {
	for _, subscriber := range eventBus.subscribers[eventType] {
		subscriber <- msg
	}
}
func (eventBus *EventBus) Subscribe(eventType string, subscriber chan<- string) {
	fmt.Println("subscribing to ", eventType)
	eventBus.subscribers[eventType] = append(eventBus.subscribers[eventType], subscriber)
}
