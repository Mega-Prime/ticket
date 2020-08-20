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
var store *ticket.Store

// create GetTicket handler
func GetTicket(w http.ResponseWriter, r *http.Request) {
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
	data, err := json.Marshal(tk)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "So sorry.")
		return
	}
	w.Write(data)
}

func ListenAndServe(s *ticket.Store) error {

	store = s

	// l := log.New(os.Stdout, "Ticket-api", log.LstdFlags)

	http.HandleFunc("/get/", GetTicket)
	ticketServer := &http.Server{
		Addr:         ":9090",
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}

	return ticketServer.ListenAndServe()

	//return http.ListenAndServe(ticketServer.Addr, sm)

}
