package main

import (
	"ticket"
	"ticket/api"
)

func main() {
	api.ListenAndServe(ticket.NewStore())
}
