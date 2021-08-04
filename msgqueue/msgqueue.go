package msgqueue

import (
	"fmt"
)

// Interface definition for events that are emitted using an EventEmitter
// Currently, the only requirement is that events are self-describing so that
// event emitter and listeners can infer an event's name.
type Event interface {
	EventName() string
}

// EventEmitter describes an interface for a class that emits events
type EventEmitter interface {
	Emit(e Event) error
}

func NewEventEmitter(mqueuetype string, mqueueconnection string, exchange string) (EventEmitter, error) {
	var emitter EventEmitter
	var err error
	switch mqueuetype {
	case "amqp":
		emitter, err = NewAMQPEventEmitter(mqueueconnection, exchange)
	case "kafka":
		emitter, err = NewKafkaEventEmitter(mqueueconnection, exchange)
	case "sqs":
		emitter, err = NewSQSEventEmitter(exchange)
	default:
		return nil, fmt.Errorf("error: Unknown MQueue driver %s", mqueuetype)
	}
	if err != nil {
		return nil, err
	}
	return emitter, nil
}

// EventListener describes an interface for a class that can listen to events.
type EventListener interface {
	Listen(events ...string) (<-chan Event, <-chan error, error)
	Mapper() EventMapper
}

func NewEventListener(mqueuetype string, mqueueconnection string, exchange string, queue string, mapper EventMapper) (EventListener, error) {
	var listener EventListener
	var err error
	switch mqueuetype {
	case "amqp":
		listener, err = NewAMQPEventListener(mqueueconnection, exchange, queue, mapper)
	case "kafka":
		listener, err = NewKafkaEventListener(mqueueconnection, exchange, mapper)
	case "sqs":
		listener, err =	NewSQSEventListener(nil, exchange, 10, 20, 10, mapper)	
	default:
		return nil, fmt.Errorf("error: Unknown MQueue driver %s", mqueuetype)
	}
	if err != nil {
		return nil, err
	}
	return listener, nil
}

type EventMapper interface {
	MapEvent(string, interface{}) (Event, error)
}
