package event

import (
	"log"
	"math/rand"
)

type Event interface {
	Handle()
	Type() string
}

type eventProvider struct {
	events []Event
}

func NewEventProvider() *eventProvider {
	return &eventProvider{
		events: []Event{},
	}
}

func (ep *eventProvider) RegisterEvent(e Event) {
	ep.events = append(ep.events, e)
}

func (ep *eventProvider) SelectEvent() {
	index := rand.Intn(len(ep.events))
	event := ep.events[index]
	log.Println(event.Type())
	event.Handle()
}
