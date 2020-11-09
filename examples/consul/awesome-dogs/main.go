package main

import (
	"net/http"
	"log"
	"encoding/json"
	"fmt"
)

type Dog struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Breed string `json:"breed"`
}

type Dogs struct {
	Dogs []Dog `json:"dogs"`
}

func LoadDogs() Dogs {
	return Dogs{
		Dogs: []Dog{
			{ID: 1, Name: "Alva", Breed: "Labrador Retriever"},
			{ID: 2, Name: "Muffin", Breed: "Labrador Retriever"},
			{ID: 3, Name: "Den", Breed: "Bob-Tailed sheep-dog"},
		},
	}
}

func main() {
	server := NewServer()
	log.Printf("Listening on %v", server.Addr)
	log.Fatal(server.ListenAndServe())
}

func NewServer() *http.Server {
	mux := http.NewServeMux()

	// Main handler that server the list of awesome dogs
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Print("Request: /")
		dogs := LoadDogs()

		data, err := json.Marshal(dogs)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "failed to encoder response: %v", err)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(data)
	})

	// A health-check endpoint used by Consul
	mux.HandleFunc("/check", func(w http.ResponseWriter, r *http.Request) {
		log.Print("Request: /check")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	server := &http.Server{
		Addr:    ":3333",
		Handler: mux,
	}

	return server
}
