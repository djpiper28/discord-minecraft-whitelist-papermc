package main

import (
	"log"
	"net/http"
)

func HealthCheckServer() {
	http.ListenAndServe(":8080", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	}))

	log.Fatal("Cannot start health check server on :8080")
}
