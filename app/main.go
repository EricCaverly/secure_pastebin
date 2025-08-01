/*********************************
 *  File     : main.go
 *  Purpose  : Backend entry point
 *  Authors  : Eric Caverly
 */

package main

import (
	"log"
	"net/http"

	"github.com/redis/go-redis/v9"
)

const (
	max_days            = 15
	max_note_size_bytes = 1024 * 30 // 30 kB
)

var rc = redis.NewClient(&redis.Options{
	Addr:     "db.spb.arpa:6379",
	Password: "",
	DB:       0,
	Protocol: 2,
})

func main() {
	// Create the HTTP server
	r := http.NewServeMux()

	r.Handle("/", http.FileServer(http.Dir("./www")))

	r.HandleFunc("GET /api/note/{id}", get_note)
	r.HandleFunc("POST /api/note", post_note)
	r.HandleFunc("GET /api/health", health_check)

	addr := ":8080"
	log.Printf("Started listening on %s\n", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalf("Web server crashed: %s\n", err.Error())
	}
}
