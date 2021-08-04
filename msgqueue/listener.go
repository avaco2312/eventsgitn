package msgqueue

import (
	"log"
)

type Handler interface{
	Handle(e interface{})
}

type EventProcessor struct {
	EventListener EventListener
	Handler Handler
	eventNames []string
}

func NewEventProcessor(mqueuetype string, mqueueconnection string, exchange string, queue string, handler Handler, mapper EventMapper, eventNames ...string) (*EventProcessor, error) {
	evlistener, err := NewEventListener(mqueuetype, mqueueconnection, exchange, queue, mapper)
	if err != nil {
		return nil, err
	}
	return &EventProcessor{EventListener: evlistener, Handler: handler, eventNames: eventNames}, nil
}

func (p *EventProcessor) ProcessEvents() {
	received, errors, err := p.EventListener.Listen(p.eventNames...)
	if err != nil {
		panic(err)
	}
	for {
		select {
		case evt := <-received:
			p.Handler.Handle(evt)
		case err = <-errors:
			log.Printf("got error while receiving event: %s\n", err)
		}
	}
}

