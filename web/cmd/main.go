package main

import "os"

func main() {
	web.ListenAndServe(os.Args[1])
}