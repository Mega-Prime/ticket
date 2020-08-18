package api

import (
	"fmt"
	"net/http"
	"ticket"
)

// create Hello handler
func hello(w http.ResponseWriter, r *http.Request) {

	s := ticket.NewStore()
	tkOpen := s.NewTicket("this is open")

	tk, err := s.GetByID(tkOpen.ID)
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}
	fmt.Fprintln(w, tk)

}

func ListenAndServe() error {

	http.HandleFunc("/", hello)

	return http.ListenAndServe(":9090", nil)

}
