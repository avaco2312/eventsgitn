package msgqueue

import (
	awsses "eventsgitn/aws"
	"encoding/json"
	"errors"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

type SQSEmitter struct {
	sqsSvc   *sqs.SQS
	QueueURL *string
}

func getOrCreateSQSQueueURL(queue *string, svc *sqs.SQS) (*string, error) {
	result, err := svc.GetQueueUrl(&sqs.GetQueueUrlInput{
		QueueName: queue,
	})
	if err == nil {
		return result.QueueUrl, nil
	}
	if v, ok := err.(awserr.Error); !ok || v.Code() != "AWS.SimpleQueueService.NonExistentQueue" {
		return nil, err
	}
	nqo, err := svc.CreateQueue(&sqs.CreateQueueInput{
		QueueName: queue,
		Attributes: map[string]*string{
			"DelaySeconds":                  aws.String("0"),
			"MessageRetentionPeriod":        aws.String("86400"),
			"ReceiveMessageWaitTimeSeconds": aws.String("20"),
		},
	})
	return nqo.QueueUrl, err
}

func NewSQSEventEmitter(queueName string) (emitter EventEmitter, err error) {
	svc := sqs.New(awsses.Sesion)
	QUrl, err := getOrCreateSQSQueueURL(aws.String(queueName), svc)
	if err != nil {
		return
	}
	emitter = &SQSEmitter{
		sqsSvc:   svc,
		QueueURL: QUrl,
	}
	return
}

func (sqsEmit *SQSEmitter) Emit(event Event) error {
	data, err := json.Marshal(event)
	if err != nil {
		return err
	}
	_, err = sqsEmit.sqsSvc.SendMessage(&sqs.SendMessageInput{
		MessageAttributes: map[string]*sqs.MessageAttributeValue{
			"event_name": {
				DataType:    aws.String("String"),
				StringValue: aws.String(event.EventName()),
			},
		},
		MessageBody: aws.String(string(data)),
		QueueUrl:    sqsEmit.QueueURL,
	})
	return err
}

type SQSListener struct {
	mapper              EventMapper
	sqsSvc              *sqs.SQS
	queueURL            *string
	maxNumberOfMessages int64
	waitTime            int64
	visibilityTimeOut   int64
}

func NewSQSEventListener(s *session.Session, queueName string, maxMsgs, wtTime, visTO int64, mapper EventMapper) (listener EventListener, err error) {
	svc := sqs.New(awsses.Sesion)
	QUrl, err := getOrCreateSQSQueueURL(aws.String(queueName), svc)
	if err != nil {
		return
	}
	listener = &SQSListener{
		sqsSvc:              svc,
		queueURL:            QUrl,
		mapper:              mapper,
		maxNumberOfMessages: maxMsgs,
		waitTime:            wtTime,
		visibilityTimeOut:   visTO,
	}
	return
}

func (sqsListener *SQSListener) Listen(events ...string) (<-chan Event, <-chan error, error) {
	if sqsListener == nil {
		return nil, nil, errors.New("SQSListener: the Listen() method was called on a nil pointer")
	}
	eventCh := make(chan Event)
	errorCh := make(chan error)
	go func() {
		for {
			sqsListener.receiveMessage(eventCh, errorCh, events...)
		}
	}()

	return eventCh, errorCh, nil
}

func (sqsListener *SQSListener) receiveMessage(eventCh chan Event, errorCh chan error, events ...string) {
	recvMsgResult, err := sqsListener.sqsSvc.ReceiveMessage(&sqs.ReceiveMessageInput{
		MessageAttributeNames: []*string{
			aws.String(sqs.QueueAttributeNameAll),
		},
		QueueUrl:            sqsListener.queueURL,
		MaxNumberOfMessages: aws.Int64(sqsListener.maxNumberOfMessages),
		WaitTimeSeconds:     aws.Int64(sqsListener.waitTime),
		VisibilityTimeout:   aws.Int64(sqsListener.visibilityTimeOut),
	})
	if err != nil {
		errorCh <- err
		return
	}
	bContinue := false
	for _, msg := range recvMsgResult.Messages {
		value, ok := msg.MessageAttributes["event_name"]
		if !ok {
			continue
		}
		eventName := aws.StringValue(value.StringValue)
		for _, event := range events {
			if strings.EqualFold(eventName, event) {
				bContinue = true
				break
			}
		}
		if !bContinue {
			continue
		}
		message := aws.StringValue(msg.Body)
		event, err := sqsListener.mapper.MapEvent(eventName, []byte(message))
		if err != nil {
			errorCh <- err
			continue
		}
		eventCh <- event
		_, err = sqsListener.sqsSvc.DeleteMessage(&sqs.DeleteMessageInput{
			QueueUrl:      sqsListener.queueURL,
			ReceiptHandle: msg.ReceiptHandle,
		})
		if err != nil {
			errorCh <- err
		}
	}
}

func (sqsListener *SQSListener) Mapper() EventMapper {
	return sqsListener.mapper
}
