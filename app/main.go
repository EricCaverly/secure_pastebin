package main

import (
	"log"
	"net/http"
	"time"

	"github.com/go-co-op/gocron/v2"
)

const (
	max_days      = 15
	trim_interval = 4 * time.Hour
)

var imdb InMemoryDB = InMemoryDB{
	notes: map[string]Note{},
}

func main() {

	s, err := gocron.NewScheduler()
	if err != nil {
		log.Fatalf("could not start gocron: %s\n", err.Error())
	}

	_, err = s.NewJob(gocron.DurationJob(trim_interval),
		gocron.NewTask(func() {
			log.Printf("starting trim job...\n")
			count := imdb.trim()
			log.Printf("trim job complete. %d notes removed\n", count)
		}))
	if err != nil {
		log.Fatalf("unable to create job: %s\n", err.Error())
	}

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
