package main

import (
	"ticket/api"
)

func main() {
	api.ListenAndServe(":8080")
}
