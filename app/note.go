package main

import "time"

type Note struct {
	Content string `json:"content"`

	AllowedIPRange string `json:"allowed_ips"`

	Created     time.Time     `json:"created"`
	ExpireAfter time.Duration `json:"expire_after"`
}
