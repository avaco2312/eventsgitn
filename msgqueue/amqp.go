package msgqueue

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/streadway/amqp"
)

type amqpEventEmitter struct {
	connection *amqp.Connection
	exchange   string
}

// NewAMQPEventEmitter creates a new event emitter.
func NewAMQPEventEmitter(url string, exchange string) (EventEmitter, error) {
	conn := <-RetryConnect(url, 5*time.Second)
	emitter := amqpEventEmitter{
		connection: conn,
		exchange:   exchange,
	}
	err := emitter.setup()
	if err != nil {
		return nil, err
	}
	return &emitter, nil
}

func (a *amqpEventEmitter) setup() error {
	channel, err := a.connection.Channel()
	if err != nil {
		return err
	}
	defer channel.Close()
	// Normally, all(many) of these options should be configurable.
	// For our example, it'll probably do.
	err = channel.ExchangeDeclare(a.exchange, "topic", true, false, false, false, nil)
	return err
}

func (a *amqpEventEmitter) Emit(event Event) error {
	channel, err := a.connection.Channel()
	if err != nil {
		return err
	}
	defer channel.Close()
	// TODO: Alternatives to JSON? Msgpack or Protobuf, maybe?
	jsonBody, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("could not JSON-serialize event: %s", err)
	}
	msg := amqp.Publishing{
		Headers:     amqp.Table{"x-event-name": event.EventName()},
		ContentType: "application/json",
		Body:        jsonBody,
	}
	err = channel.Publish(a.exchange, event.EventName(), false, false, msg)
	return err
}

const eventNameHeader = "x-event-name"

type amqpEventListener struct {
	connection *amqp.Connection
	exchange   string
	queue      string
	mapper     EventMapper
}

// NewAMQPEventListener creates a new event listener.
func NewAMQPEventListener(url string, exchange string, queue string, mapper EventMapper) (EventListener, error) {

	//(conn *amqp.Connection, exchange string, queue string) (EventListener, error) {

	conn := <-RetryConnect(url, 5*time.Second)
	listener := amqpEventListener{
		connection: conn,
		exchange:   exchange,
		queue:      queue,
		mapper:     mapper,
	}
	err := listener.setup()
	if err != nil {
		return nil, err
	}
	return &listener, nil
}

// NewAMQPEventListenerFromEnvironment will create a new event listener from

func (a *amqpEventListener) setup() error {
	channel, err := a.connection.Channel()
	if err != nil {
		return err
	}
	defer channel.Close()
	err = channel.ExchangeDeclare(a.exchange, "topic", true, false, false, false, nil)
	if err != nil {
		return err
	}
	_, err = channel.QueueDeclare(a.queue, true, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("could not declare queue %s: %s", a.queue, err)
	}
	return nil
}

// Listen configures the event listener to listen for a set of events that are
// specified by name as parameter.
// This method will return two channels: One will contain successfully decoded
// events, the other will contain errors for messages that could not be
// successfully decoded.
func (l *amqpEventListener) Listen(eventNames ...string) (<-chan Event, <-chan error, error) {
	channel, err := l.connection.Channel()
	if err != nil {
		return nil, nil, err
	}
	// Create binding between queue and exchange for each listened event type
	for _, event := range eventNames {
		if err := channel.QueueBind(l.queue, event, l.exchange, false, nil); err != nil {
			return nil, nil, fmt.Errorf("could not bind event %s to queue %s: %s", event, l.queue, err)
		}
	}
	msgs, err := channel.Consume(l.queue, "", false, false, false, false, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("could not consume queue: %s", err)
	}
	events := make(chan Event)
	errors := make(chan error)
	go func() {
		for msg := range msgs {
			rawEventName, ok := msg.Headers[eventNameHeader]
			if !ok {
				errors <- fmt.Errorf("message did not contain %s header", eventNameHeader)
				msg.Nack(false, false)
				continue
			}
			eventName, ok := rawEventName.(string)
			if !ok {
				errors <- fmt.Errorf("header %s did not contain string", eventNameHeader)
				msg.Nack(false, false)
				continue
			}
			event, err := l.mapper.MapEvent(eventName, msg.Body)
			if err != nil {
				errors <- fmt.Errorf("could not unmarshal event %s: %s", eventName, err)
				msg.Nack(false, false)
				continue
			}
			msg.Ack(false)
			events <- event
		}
	}()
	return events, errors, nil
}

func (l *amqpEventListener) Mapper() EventMapper {
	return l.mapper
}

// RetryConnect implements a retry mechanism for establishing the AMQP connection.
// This is necessary in container environments where individual components may be
// started out-of-order, so we might have to wait for upstream services like RabbitMQ
// to actually become available.
//
// Alternatives:
//   - use an entrypoint script in your container in which you wait for the service
//     to be available [1]
//   - use a script like wait-for-it [2] in your entrypoint
//
// [1] http://stackoverflow.com/q/25503412/1995300
// [1] https://github.com/vishnubob/wait-for-it/blob/master/wait-for-it.sh
func RetryConnect(amqpURL string, retryInterval time.Duration) chan *amqp.Connection {
	result := make(chan *amqp.Connection)
	go func() {
		defer close(result)
		for {
			conn, err := amqp.Dial(amqpURL)
			if err == nil {
				result <- conn
				return
			}
			time.Sleep(retryInterval)
		}
	}()
	return result
}