package rest

import (
	"encoding/hex"
	"fmt"

	"golang.org/x/net/context"

	"eventsgitn/contractsp"
	"eventsgitn/bookservice/store"
)

type grpchandler struct {
	store   store.Store
	contractsp.UnimplementedEventsServer
}

func NewGRPCHandler(istore store.Store) (*grpchandler, error) {
	return &grpchandler{store: istore}, nil
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
