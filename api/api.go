package api

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"path"
	"ticket"
	"time"
)

// bad idea! fix this
// problem: how do we access the store from within the handler?
var (
	store  ticket.Store
	Logger *log.Logger
)

func ListenAndServe(addr, dbURI, dbName, collection string) error {
	var err error
	store, err = ticket.NewMongoStore(context.Background(), dbURI, dbName, collection)
	if err != nil {
		return err
	}
	log.Println("Did we get here?")
	Logger = log.New(os.Stderr, "", log.LstdFlags)

	Logger.Println("Server started on: ", addr)
	sm := http.NewServeMux()
	sm.HandleFunc("/get/", GetTicket)
	sm.HandleFunc("/create", createTicket)
	sm.HandleFunc("/all", GetAllTickets)
	sm.HandleFunc("/healthz", healthz)

	ticketServer := &http.Server{
		Handler:      sm,
		Addr:         addr,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}

	return ticketServer.ListenAndServe()

	// return http.ListenAndServe(ticketServer.Addr, sm)
}

//  GetTicket handler
func GetTicket(response http.ResponseWriter, request *http.Request) {
	Logger.Println(request.Method, request.URL)
	_, ID := path.Split(request.URL.Path)

	tk, err := store.GetByID(ticket.ID(ID))
	if err != nil {
		response.WriteHeader(http.StatusNotFound)
		response.Write([]byte(`{"message": "` + err.Error() + `"}`))
		//fmt.Fprintln(response, err)
		return
	}
	err = json.NewEncoder(response).Encode(tk)
	if err != nil {
		// log.Println(err)
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"message": "` + err.Error() + `"}`))
		//fmt.Fprintln(response, "So sorry.")
		return
	}
}

//  write get all tickets func:
func GetAllTickets(response http.ResponseWriter, request *http.Request) {
	Logger.Println(request.Method, request.URL)
	path.Split(request.URL.Path)
	tks, err := store.GetAll()
	if err != nil {
		response.WriteHeader(http.StatusNotFound)
		response.Write([]byte(`{"message": "` + err.Error() + `"}`))
		return
	}

	err = json.NewEncoder(response).Encode(tks)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"message": "` + err.Error() + `"}`))
		return
	}

}

// healthz check 200:
func healthz(w http.ResponseWriter, r *http.Request) {
	Logger.Println(r.Method, r.URL)
}

// CreateTicket Handler
func createTicket(w http.ResponseWriter, r *http.Request) {
	Logger.Println(r.Method, r.URL)

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
