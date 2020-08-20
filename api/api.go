package api

import (
	"fmt"
	"net/http"
	"path"
	"strconv"
	"ticket"
)

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
	fmt.Fprintln(w, tk)
}

func ListenAndServe(s *ticket.Store) error {

	store = s

	// l := log.New(os.Stdout, "Ticket-api", log.LstdFlags)

	// hh := Hello(l)

	// sm := http.NewServeMux()
	// sm.Handle("/", hh)
	// sm.Handle("/1", hh)

	// ticketServer := &http.Server{
	// 	Addr:         ":9090",
	// 	Handler:      sm,
	// 	IdleTimeout:  120 * time.Second,
	// 	ReadTimeout:  1 * time.Second,
	// 	WriteTimeout: 1 * time.Second,
	// }

	// go func() {
	// 	err := ticketServer.ListenAndServe()
	// 	if err != nil {
	// 		l.Fatal(err)
	// 	}
	// }()

	http.HandleFunc("/get/", GetTicket)

	return http.ListenAndServe(":9090", nil)

	//return http.ListenAndServe(ticketServer.Addr, sm)

}
