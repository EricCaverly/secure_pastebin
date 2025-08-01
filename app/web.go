/*********************************
 *  File     : web.go
 *  Purpose  : Backend web logic for API endpoints
 *  Authors  : Eric Caverly
 */

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Resp struct {
	Success bool   `json:"success"`
	Data    any    `json:"data"`
	Message string `json:"message"`
}

func write_error(w http.ResponseWriter, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	bdy := Resp{
		Success: false,
		Message: msg,
		Data:    nil,
	}

	d, _ := json.Marshal(bdy)

	w.Write(d)
}

func write_success(w http.ResponseWriter, msg string, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	bdy := Resp{
		Success: true,
		Message: msg,
		Data:    data,
	}

	d_buff, err := json.Marshal(bdy)
	if err != nil {
		write_error(w, "Failed to marshal data")
		log.Printf("Failed to format a success body: %s\n", err.Error())
		return
	}

	w.Write(d_buff)
}

func get_note(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	log.Printf("get_note called from %s on uuid:%s\n", r.RemoteAddr, id)

	n, err := rc.HGetAll(context.Background(), id).Result()
	if err != nil {
		write_error(w, fmt.Sprintf("Failed to get note: %s", err.Error()))
		return
	}

	// If using a reverse proxy, which is assumed... Otherwise use r.RemoteAddr
	remote_addr := r.Header.Get("X-Real-Ip")
	if remote_addr == "" {
		remote_addr = strings.Split(r.RemoteAddr, ":")[0]
	}
	log.Printf("IP from Traefik: %s\n", remote_addr)

	allowed_ips, ok := n["allowed_ips"]
	if !ok {
		write_error(w, "Note not found")
		return
	}

	allowed, err := within_ranges(remote_addr, allowed_ips)
	if err != nil {
		write_error(w, "Failed to check if IP was valid")
		return
	}

	if !allowed {
		write_error(w, "You are not allowed to access this note! (IP address forbidden)")
		return
	}

	write_success(w, "Found note", n)
}

func post_note(w http.ResponseWriter, r *http.Request) {
	log.Printf("post_note called from %s\n", r.RemoteAddr)

	r.ParseForm()

	content, ok := r.Form["content"]
	if !ok {
		write_error(w, "Missing content in request")
		return
	}

	if len(content[0]) > max_note_size_bytes {
		write_error(w, fmt.Sprintf("Note too large! Max size: %d bytes", max_note_size_bytes))
		return
	}

	allowed_ips, ok := r.Form["allowed_ips"]
	if !ok {
		write_error(w, "Missing allowed_ips in request")
		return
	}

	if err := check_valid_ranges(allowed_ips[0]); err != nil {
		write_error(w, fmt.Sprintf("Invalid IP range (%s). Please enter as 1.1.1.0/24, 2.2.0.0/16", err.Error()))
		return
	}

	dte, ok := r.Form["days_until_expire"]
	if !ok {
		write_error(w, "Missing days_until_expire in request")
		return
	}

	ndays, err := strconv.Atoi(dte[0])
	if err != nil || ndays < 0 || ndays > max_days {
		write_error(w, "Invalid number of days")
		return
	}

	n := Note{
		Content:        content[0],
		AllowedIPRange: allowed_ips[0],
	}

	id := uuid.NewString()

	rc.HSet(context.Background(), id, n)
	rc.Expire(context.Background(), id, time.Duration(ndays)*24*time.Hour)

	write_success(w, "Note created!", id)
}

func health_check(w http.ResponseWriter, _ *http.Request) {
	write_success(w, "Alive", nil)
}
