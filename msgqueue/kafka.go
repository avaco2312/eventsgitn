package msgqueue

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/Shopify/sarama"
)

type kafkaEventEmitter struct {
	producer sarama.SyncProducer
	topic    string
}

func NewKafkaEventEmitter(url string, topic string) (EventEmitter, error) {
	brokers := []string{url}
	client := <-RetryConnectKafka(brokers, 5*time.Second)
	client.Config().Producer.Return.Successes = true
	producer, err := sarama.NewSyncProducerFromClient(client)
	if err != nil {
		return nil, err
	}
	emitter := kafkaEventEmitter{
		producer: producer,
		topic:    topic,
	}
	return &emitter, nil
}

func (k *kafkaEventEmitter) Emit(evt Event) error {
	jsonBody, err := json.Marshal(evt)
	if err != nil {
		return err
	}
	msg := &sarama.ProducerMessage{
		Topic: k.topic,
		Key:   sarama.ByteEncoder(evt.EventName()),
		Value: sarama.ByteEncoder(jsonBody),
	}
	_, _, err = k.producer.SendMessage(msg)
	if err != nil {
		log.Println(err.Error())
	}
	return err
}

type kafkaEventListener struct {
	consumer   sarama.Consumer
	partitions []int32
	mapper     EventMapper
	topic      string
}

func NewKafkaEventListener(url string, topic string, mapper EventMapper) (EventListener, error) {
	brokers := []string{url}
	partitions := []int32{}
	client := <-RetryConnectKafka(brokers, 5*time.Second)
	consumer, err := sarama.NewConsumerFromClient(client)
	if err != nil {
		return nil, err
	}
	listener := &kafkaEventListener{
		consumer:   consumer,
		partitions: partitions,
		mapper:     mapper,
		topic:      topic,
	}
	return listener, nil
}

func (k *kafkaEventListener) Listen(events ...string) (<-chan Event, <-chan error, error) {
	var err error
	topic := k.topic
	results := make(chan Event)
	errors := make(chan error)
	partitions := k.partitions
	if len(partitions) == 0 {
		partitions, err = k.consumer.Partitions(topic)
		if err != nil {
			return nil, nil, err
		}
	}
	for _, partition := range partitions {
		pConsumer, err := k.consumer.ConsumePartition(topic, partition, 0)
		if err != nil {
			return nil, nil, err
		}
		go func() {
			for msg := range pConsumer.Messages() {
				event, err := k.mapper.MapEvent(string(msg.Key), msg.Value)
				if err != nil {
					errors <- fmt.Errorf("could not map message: %v", err)
					continue
				}
				results <- event
			}
		}()
		go func() {
			for err := range pConsumer.Errors() {
				errors <- err
			}
		}()
	}
	return results, errors, nil
}

func (k *kafkaEventListener) Mapper() EventMapper {
	return k.mapper
}

// RetryConnect implements a retry mechanism for establishing the Kafka connection.
// This is necessary in container environments where individual components may be
// started out-of-order, so we might have to wait for upstream services like Kafka
// to actually become available.
//
// Alternatives:
//   - use an entrypoint script in your container in which you wait for the service
//     to be available [1]
//   - use a script like wait-for-it [2] in your entrypoint
//
// [1] http://stackoverflow.com/q/25503412/1995300
// [1] https://github.com/vishnubob/wait-for-it/blob/master/wait-for-it.sh
func RetryConnectKafka(brokers []string, retryInterval time.Duration) chan sarama.Client {
	result := make(chan sarama.Client)
	go func() {
		defer close(result)
		for {
			config := sarama.NewConfig()
			conn, err := sarama.NewClient(brokers, config)
			if err == nil {
				result <- conn
				return
			}
			log.Printf("Kafka connection failed with error (retrying in %s): %s", retryInterval.String(), err)
			time.Sleep(retryInterval)
		}
	}()
	return result
}
