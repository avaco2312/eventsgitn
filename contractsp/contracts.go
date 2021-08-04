package contractsp

type EventCreated struct {
	Event *Event
}

// EventName returns the event's name
func (e EventCreated) EventName() string {
	return "event.Created"
}