package handlers

import (
	"log"
	"encoding/json"
	"net/http"
)

func respondWithError(w http.ResponseWriter, code int, msg string) {
	if code > 499 {
		log.Println("Responding with 5XX error: ", msg)
	}

	type errorResponse struct {
		Error string `json:"error"`
	}

	respondWithJson(w, code, errorResponse{
		Error: msg,
	})
}

func respondWithJson(w http.ResponseWriter, code int, payload interface{}) {
	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Build to marshal json response: %v", payload)
		w.WriteHeader(500)
		return
	}

	w.Header().Add("Content-Type:", "application/json")
	w.WriteHeader(code)
	w.Write(data)
}