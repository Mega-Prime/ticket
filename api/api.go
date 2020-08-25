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

func ListenAndServe(s *ticket.Store) error {
	log.Println("Server started")
	store = s

	http.HandleFunc("/get/", GetTicket)
	http.HandleFunc("/create", createTicket)

	ticketServer := &http.Server{
		Addr:         ":9090",
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}

	return ticketServer.ListenAndServe()

	//return http.ListenAndServe(ticketServer.Addr, sm)

}

// create GetTicket handler
func GetTicket(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
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
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "So sorry.")
		return
	}
}

// CreateTicket Handler
func createTicket(w http.ResponseWriter, r *http.Request) {
	log.Println("handle POST request")
	log.Println(r.URL)
	tk := ticket.Ticket{}
	// decoder := json.NewDecoder(r.Body)
	// error := decoder.Decode(&tk)
	// if error != nil {
	// 	log.Println(error.Error())
	// 	http.Error(w, "Unable to unmarshal json data", http.StatusBadRequest)
	// 	return
	// }
	err := tk.FromJson(r.Body)
	if err != nil {
		http.Error(w, "unable to unmarshal json data", http.StatusBadRequest)

	}
	ticket.CreateTicket(tk)

}
