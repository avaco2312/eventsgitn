package store

import (
	"encoding/hex"
	"eventsgitn/contractsp"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type mongoStore struct {
	session *mgo.Session
	db      string
}

func NewMongoStore(conn string, ndb string) (*mongoStore, error) {
	ses, err := mgo.Dial(conn)
	if err != nil {
		return nil, err
	}
	return &mongoStore{session: ses, db: ndb}, nil
}

func (m *mongoStore) SearchId(id interface{}) (*contractsp.Event, error) {
	bid := hex.EncodeToString(id.([]byte))
	ses := m.session.Copy()
	defer ses.Close()
	var event contractsp.Event
	err := ses.DB(m.db).C("events").Find(bson.M{"_id": bid}).One(&event)
	if err != nil {
		if err == mgo.ErrNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &event, nil
}

func (m *mongoStore) SearchName(nEvent string) (*contractsp.Event, error) {
	ses := m.session.Copy()
	defer ses.Close()
	var event contractsp.Event
	err := ses.DB(m.db).C("events").Find(bson.M{"name": nEvent}).One(&event)
	if err != nil {
		if err == mgo.ErrNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &event, nil
}

func (m *mongoStore) SearchAll() (contractsp.ArrayEvent, error) {
	ses := m.session.Copy()
	defer ses.Close()
	var events contractsp.ArrayEvent
	err := ses.DB(m.db).C("events").Find(nil).All(&events.Events)
	if err != nil {
		return events, err
	}
	return events, nil
}

func (m *mongoStore) AddEvent(ev *contractsp.Event) error {
	ses := m.session.Copy()
	defer ses.Close()
	err := ses.DB(m.db).C("events").Insert(ev)
	return err
}
