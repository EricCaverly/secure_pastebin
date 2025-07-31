/*********************************
 *  File     : imbd.go
 *  Purpose  : In Memory DataBase
 *             Stores notes within RAM. If expansion is required in the future, replace with redis
 *  Authors  : Eric Caverly
 */

package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
)

type InMemoryDB struct {
	mtx         sync.Mutex
	notes       map[string]Note
	total_count int
}

func (db *InMemoryDB) trim() int {
	db.mtx.Lock()
	defer db.mtx.Unlock()

	count := 0

	for k, v := range db.notes {
		if time.Now().After((v.Created.Add(v.ExpireAfter))) {
			delete(db.notes, k)
			count++
		}
	}

	return count
}

func (db *InMemoryDB) push(note Note) (string, error) {
	db.mtx.Lock()
	defer db.mtx.Unlock()

	if db.total_count > max_notes {
		return "", fmt.Errorf("max amount of notes reached (%d)", max_notes)
	}

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
