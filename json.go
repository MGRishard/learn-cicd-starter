package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func respondWithError(w http.ResponseWriter, code int, msg string, logErr error) {
	if logErr != nil {
		log.Println(logErr)
	}
	if code > 499 {
		log.Printf("Responding with 5XX error: %s", msg)
	}
	type errorResponse struct {
		Error string `json:"error"`
	}
	respondWithJSON(w, code, errorResponse{
		Error: msg,
	})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	dat, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}
	// 1. Fix G705 (XSS): Explicitly set Content-Type so the browser
	// doesn't try to "guess" if the data is executable HTML.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	// 2. Fix G104 (Unhandled Error): Capture the error return.
	// Also remove the duplicate w.Write(dat) line if you have two.
	_, err = w.Write(dat)
	if err != nil {
		log.Printf("Failed to write response: %v", err)
	}
}
