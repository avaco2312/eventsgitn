package rest

import (
	"log"
	"net"
	"net/http"

	"eventsgitn/contractsp"
	"eventsgitn/eventservice/store"
	"eventsgitn/msgqueue"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func ServeApi(store store.Store, emitter msgqueue.EventEmitter, endpoint string, path string) chan error {
	go func() {
		h := mux.NewRouter()
		h.Handle("/metrics", promhttp.Handler())
		http.ListenAndServe(":8072", handlers.CORS()(h))
	}()
	handler, _ := NewHandler(store, emitter)
	grpchandler, _ := NewGRPCHandler(store, emitter)
	r := mux.NewRouter()
	eventsRouter := r.PathPrefix("/" + path).Subrouter()
	eventsRouter.Methods("GET").Path("/{search}/{params}").HandlerFunc(handler.search)
	eventsRouter.Methods("GET").Path("").HandlerFunc(handler.searchAll)
	eventsRouter.Methods("POST").Path("").HandlerFunc(handler.addEvent)
	listener, err := net.Listen("tcp", ":8091")
	if err != nil {
		log.Fatal("Error de listener TCP")
	}
	s := grpc.NewServer()
	contractsp.RegisterEventsServer(s, grpchandler)
	reflection.Register(s)
	cherr := make(chan error)
	go func() {
		cherr <- http.ListenAndServe(endpoint, handlers.CORS()(r))
		//cherr <- http.ListenAndServeTLS(endpoint, "cert.pem", "key.pem", r)
	}()
	go func() {
		cherr <- s.Serve(listener)
	}()
	return cherr
}
