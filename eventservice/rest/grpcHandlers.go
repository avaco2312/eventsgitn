package rest

import (
	"encoding/hex"
	"fmt"
	"time"

	"golang.org/x/net/context"

	"eventsgitn/contractsp"
	"eventsgitn/eventservice/store"
	"eventsgitn/msgqueue"
)

type grpchandler struct {
	store   store.Store
	emitter msgqueue.EventEmitter
	contractsp.UnimplementedEventsServer
}

func NewGRPCHandler(istore store.Store, iemitter msgqueue.EventEmitter) (*grpchandler, error) {
	return &grpchandler{store: istore, emitter: iemitter}, nil
}

func (g *grpchandler) SearchId(ctx context.Context, in *contractsp.Id) (*contractsp.Event, error) {
	id, err := hex.DecodeString(in.Id)
	if err != nil || len(id) != 12 {
		return nil, fmt.Errorf("bad hexadecimal id, you can search by id via /id/(12 bytes hex)")
	}
	event, err := g.store.SearchId(id)
	if event == nil {
		return nil, fmt.Errorf("not found")
	}
	if err != nil {
		return nil, fmt.Errorf("error del servidor")
	}
	return event, nil
}

func (g *grpchandler) SearchName(ctx context.Context, in *contractsp.Name) (*contractsp.Event, error) {
	event, err := g.store.SearchName(in.Name)
	if event == nil {
		return nil, fmt.Errorf("not found")
	}
	if err != nil {
		return nil, fmt.Errorf("error del servidor")
	}
	return event, nil
}

func (g *grpchandler) SearchAll(ctx context.Context, in *contractsp.Empty) (*contractsp.ArrayEvent, error) {
	events, err := g.store.SearchAll()
	if err != nil {
		return nil, fmt.Errorf("error del servidor")
	}
	if len(events.Events) == 0 {
		return nil, fmt.Errorf("not found")
	}
	return &events, nil
}

func (g *grpchandler) AddEvent(ctx context.Context, in *contractsp.Event) (*contractsp.Id, error) {
	now := time.Now()
	id, err := g.store.AddEvent(in)
	if err != nil {
		return nil, fmt.Errorf("error del servidor")
	}
	g.emitter.Emit(contractsp.EventCreated{Event: in})
	cid := contractsp.Id{
		Id: fmt.Sprintf("%x", id.(string)),
	}
	eventsAddDelay.Observe(float64(time.Since(now)) / 1e6)
	return &cid, nil
}
