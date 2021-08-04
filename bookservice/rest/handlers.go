package rest

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"eventsgitn/bookservice/store"

	"github.com/gorilla/mux"
)

type handler struct {
	store store.Store
}

func NewHandler(istore store.Store) (*handler, error) {
	return &handler{store: istore}, nil
}

func (h *handler) search(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json;charset=utf8")
	vars := mux.Vars(r)
	criterio, ok := vars["search"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, `{"error": "No search criteria found, you can either search by id via /id/(12 bytes hex) or search by name via /name/coldplayconcert"}`)
		return
	}
	param, ok := vars["params"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, `{"error": "No search keys found, you can either search by id via /id/(12 bytes hex) or search by name via /name/coldplayconcert"}`)
		return
	}
	switch strings.ToLower(criterio) {
	case "id":
		id, err := hex.DecodeString(param)
		if err != nil || len(id) != 12 {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, `{"error": "Bad hexadecimal id, you can search by id via /id/(12 bytes hex)"}`)
			return
		}
		event, err := h.store.SearchId(id)
		if event == nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(&event)
	case "name":
		event, err := h.store.SearchName(param)
		if event == nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(&event)
	default:
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, `{"error": "Bad search criteria found, you can either search by id via /id/(12 bytes hex) or search by name via /name/coldplayconcert"}`)
		return
	}
}

func (h *handler) searchAll(w http.ResponseWriter, r *http.Request) {
	events, err := h.store.SearchAll()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if len(events.Events) == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json;charset=utf8")
	json.NewEncoder(w).Encode(&events)
}
