package api

import (
	"fmt"
	"net/http"
	"ticket"
)

var store *ticket.Store

// create Hello handler
func Hello(w http.ResponseWriter, r *http.Request) {
	storedID := store.NewTicket("New test ticket")
	ID := storedID.ID

	tk, err := store.GetByID(ID)
	if err != nil {
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

	http.HandleFunc("/", Hello)
	http.HandleFunc("/1", Hello)

	return http.ListenAndServe(":9090", nil)

	//return http.ListenAndServe(ticketServer.Addr, sm)

}
