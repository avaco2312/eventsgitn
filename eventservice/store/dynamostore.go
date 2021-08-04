package store

import (
	"encoding/hex"
	"encoding/json"

	awsses "eventsgitn/aws"
	"eventsgitn/contractsp"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"gopkg.in/mgo.v2/bson"
)

type eventAWS struct {
	ID   string
	Name string `dynamodbav:"EventName"`
	Body string
}

type DynamoStore struct {
	service *dynamodb.DynamoDB
	db      string
}

func NewDynamoStore(db string) (*DynamoStore, error) {
	return &DynamoStore{
		service: dynamodb.New(awsses.Sesion),
		db:      db,
	}, nil
}

func (d *DynamoStore) AddEvent(event *contractsp.Event) (interface{}, error) {
	bid := bson.NewObjectId()
	event.Id = hex.EncodeToString([]byte(bid))
	event.Location.Id = hex.EncodeToString([]byte(bson.NewObjectId()))
	evaws := eventAWS{}
	evaws.ID = event.Id
	evaws.Name = event.Name
	sbody, err := json.Marshal(event)
	if err != nil {
		return nil, err
	}
	evaws.Body = string(sbody)
	av, err := dynamodbattribute.MarshalMap(evaws)
	if err != nil {
		return nil, err
	}
	_, err = d.service.PutItem(&dynamodb.PutItemInput{
		TableName: aws.String(d.db),
		Item:      av,
	})
	if err != nil {
		return nil, err
	}
	return string(bid), nil
}

func (d *DynamoStore) SearchId(id interface{}) (*contractsp.Event, error) {
	bid := hex.EncodeToString([]byte(id.([]byte)))
	input := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"ID": {
				S: aws.String(bid),
			},
		},
		TableName: aws.String(d.db),
	}
	result, err := d.service.GetItem(input)
	if err != nil {
		return nil, err
	}
	evaws := eventAWS{}
	err = dynamodbattribute.UnmarshalMap(result.Item, &evaws)
	if err != nil {
		return nil, err
	}
	var event contractsp.Event
	err = json.Unmarshal([]byte(evaws.Body), &event)
	if err != nil {
		return nil, err
	}
	return &event, nil
}

func (d *DynamoStore) SearchName(name string) (*contractsp.Event, error) {
	input := &dynamodb.QueryInput{
		KeyConditionExpression: aws.String("EventName = :n"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":n": {
				S: aws.String(name),
			},
		},
		IndexName: aws.String("EventName-index"),
		TableName: aws.String(d.db),
	}
	result, err := d.service.Query(input)
	if err != nil {
		return nil, err
	}
	evaws := eventAWS{}
	var event contractsp.Event
	if len(result.Items) > 0 {
		err = dynamodbattribute.UnmarshalMap(result.Items[0], &evaws)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal([]byte(evaws.Body), &event)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, nil
	}
	return &event, err
}

func (d *DynamoStore) SearchAll() (contractsp.ArrayEvent, error) {
	var events contractsp.ArrayEvent
	input := &dynamodb.ScanInput{
		TableName: aws.String(d.db),
	}
	result, err := d.service.Scan(input)
	if err != nil {
		return events, err
	}
	var evawss []eventAWS
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &evawss)
	if err != nil {
		return events, err
	}
	var event contractsp.Event
	for _, evaws := range evawss {
		err = json.Unmarshal([]byte(evaws.Body), &event)
		if err != nil {
			return events, err
		}
		events.Events = append(events.Events, &event)
	}
	return events, err
}
