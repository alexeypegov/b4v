package model

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/boltdb/bolt"
)

const (
	notesBucket = "notes"
)

var (
	uuidRegexp = regexp.MustCompile("([\\s\\pP\\pS]+)")
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

// GetNote loads Note by its uuid
func GetNote(uuid string, db *DB) (*Note, error) {
	bytes, err := db.Get(notesBucket, uuid)
	if err != nil {
		return nil, err
	}

	result := new(Note)
	json.Unmarshal(bytes, &result)
	return result, nil
}

func genUUID(title string) string {
	time := time.Now()

	replaced := strings.ToLower(uuidRegexp.ReplaceAllLiteralString(title, "-"))
	return fmt.Sprintf("%4d%2d%2d-%s", time.Year(), time.Month(), time.Day(), replaced)
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

		if len(note.UUID) == 0 {
			note.UUID = genUUID(note.Title)
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
