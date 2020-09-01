package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"path"
	"strconv"
	"ticket"
	"time"
)

// bad idea! fix this
// problem: how do we access the store from within the handler?

var store *ticket.Store

func ListenAndServe(addr string) error {
	log.Println("Server started on: ", addr)
	store = ticket.NewStore()

	sm := http.NewServeMux()
	sm.HandleFunc("/get/", GetTicket)
	sm.HandleFunc("/create", createTicket)
	sm.HandleFunc("/healthz", healthz)

	ticketServer := &http.Server{
		Handler:      sm,
		Addr:         addr,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}

	return ticketServer.ListenAndServe()

	//return http.ListenAndServe(ticketServer.Addr, sm)

}

// create GetTicket handler
func GetTicket(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Method, r.URL)
	_, rawID := path.Split(r.URL.Path)
	ID, err := strconv.Atoi(rawID)
	if err != nil {
		w.WriteHeader(http.StatusTeapot)
		fmt.Fprintf(w, "Invalid ticket ID %q", rawID)
		return
	}
	tk, err := store.GetByID(ID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w, err)
		return
	}
	err = json.NewEncoder(w).Encode(tk)
	if err != nil {
		//log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "So sorry.")
		return
	}
}

//healthz check 200:
func healthz(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Method, r.URL)
}

// CreateTicket Handler
func createTicket(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Method, r.URL)

	tk := ticket.Ticket{}

	err := tk.FromJSON(r.Body)
	if err != nil {
		http.Error(w, "unable to unmarshal json data", http.StatusBadRequest)
		return
	}
	ID, err := store.AddTicket(tk)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	retrieve, err := store.GetByID(ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusCreated)
	retrieve.ToJSON(w)

}
