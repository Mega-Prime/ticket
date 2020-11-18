package web

import (
	"net/http"
	"time"
)

func ListenAndServeWithTLS(addr, certFile, keyFile string) error {
	sm := http.NewServeMux()
	sm.HandleFunc("/healthz", healthz)

	w := &http.Server{
		Handler:      sm,
		Addr:         addr,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}
	return w.ListenAndServeTLS(certFile, keyFile)
}

// healthz check 200:
func healthz(w http.ResponseWriter, r *http.Request) {
}
