package main

import (
	"fmt"
	"sync"

	"github.com/google/uuid"
)

type InMemoryDB struct {
	mtx         sync.Mutex
	notes       map[string]Note
	total_count int
}

func (db *InMemoryDB) push(note Note) (string, error) {
	db.mtx.Lock()
	defer db.mtx.Unlock()

	id := uuid.NewString()
	db.notes[id] = note

	db.total_count++

	return id, nil
}

func (db *InMemoryDB) fetch(id string) (Note, error) {
	db.mtx.Lock()
	defer db.mtx.Unlock()

	n, ok := db.notes[id]
	if !ok {
		return Note{}, fmt.Errorf("id not found, note does not exist")
	}

	return n, nil
}

func (db *InMemoryDB) pop(id string) (Note, error) {
	db.mtx.Lock()
	defer db.mtx.Unlock()

	n, ok := db.notes[id]
	if !ok {
		return Note{}, fmt.Errorf("id not found, note does not exist")
	}

	delete(db.notes, id)
	db.total_count--

	return n, nil
}
