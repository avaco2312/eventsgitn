package rest

import (
	"log"
	"net"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"eventsgitn/bookservice/store"
	"eventsgitn/contractsp"
)

func ServeApi(store store.Store, endpoint string, path string) chan error {
	go func() {
		h := mux.NewRouter()
		h.Handle("/metrics", promhttp.Handler())
		http.ListenAndServe(":8073", handlers.CORS()(h))
	}()
	handler, _ := NewHandler(store)
	grpchandler, _ := NewGRPCHandler(store)
	r := mux.NewRouter()
	eventsRouter := r.PathPrefix("/" + path).Subrouter()
	eventsRouter.Methods("GET").Path("/events/{search}/{params}").HandlerFunc(handler.search)
	eventsRouter.Methods("GET").Path("/events").HandlerFunc(handler.searchAll)
	listener, err := net.Listen("tcp", ":8092")
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
