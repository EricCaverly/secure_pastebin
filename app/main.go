package main

import (
	"log"
	"net/http"
)

var imdb InMemoryDB = InMemoryDB{
	notes: map[string]Note{},
}

func main() {
	r := http.NewServeMux()

	r.Handle("/", http.FileServer(http.Dir("./www")))

	r.HandleFunc("GET /api/note/{id}", get_note)
	r.HandleFunc("POST /api/note", post_note)

	addr := ":8080"
	log.Printf("Started listening on %s\n", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalf("Web server crashed: %s\n", err.Error())
	}

}
