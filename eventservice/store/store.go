package store

import (
	"eventsgitn/contractsp"
	"fmt"
)

type Store interface {
	SearchId(interface{}) (*contractsp.Event, error)
	SearchName(string) (*contractsp.Event, error)
	SearchAll() (contractsp.ArrayEvent, error)
	AddEvent(*contractsp.Event) (interface{}, error)
}

func NewStore(dbType string, connString string, db string) (Store, error) {
	var st Store
	var err error
	switch dbType {
	case "mongo":
		st, err = NewMongoStore(connString, db)
	case "dynamo":
			st, err = NewDynamoStore(db)
	default:
		return nil, fmt.Errorf("error: Unknown DB driver %s", dbType)
	}
	if err != nil {
		return nil, err
	}
	return st, nil
}
