package model

import (
	"encoding/json"
	"time"

	"github.com/boltdb/bolt"
)

const (
	notesBucket = "notes"
)

const (
	// Draft determines whatever this note is published or not
	Draft byte = 1 << iota

	// PlainHTML format of the entry
	PlainHTML
)

// Note is a note, yeah
type Note struct {
	UUID      string    `json:"uuid"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Tags      []string  `json:"tags"`
	CreatedAt time.Time `json:"created_at"`
	Flags     byte      `json:"flags"`
}

// SaveAll save all the given Notes
func SaveAll(notes []Note, db *DB) error {
	for _, n := range notes {
		if err := n.Save(false, db); err != nil {
			return err
		}
	}

	return nil
}

// Load loads Note by its uuid
func (note *Note) Load(uuid string, db *DB) error {
	bytes, err := db.Get(notesBucket, uuid)
	if err != nil {
		return err
	}

	json.Unmarshal(bytes, &note)
	return nil
}

// Save Note to a storage
func (note *Note) Save(draft bool, db *DB) error {
	note.Flags |= Draft
	note.CreatedAt = time.Now()

	if err := db.Handle.Update(func(tx *bolt.Tx) error {
		bucketNotes, err := tx.CreateBucketIfNotExists([]byte(notesBucket))
		if err != nil {
			return err
		}

		if len(note.UUID) < 1 {
			note.UUID = "todo" // todo
		}

		jsonNote, _ := json.Marshal(note)
		if err := bucketNotes.Put([]byte(note.UUID), []byte(jsonNote)); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}
